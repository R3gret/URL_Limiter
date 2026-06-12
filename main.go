package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// --- RATE LIMITER SETTINGS ---
const (
	RateLimitRequests = 5
	RateLimitWindow   = 10 * time.Second
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	i := &RateLimiter{
		visitors: make(map[string]*visitor),
	}
	go func() {
		for {
			time.Sleep(time.Minute)
			i.mu.Lock()
			for id, v := range i.visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(i.visitors, id)
				}
			}
			i.mu.Unlock()
		}
	}()
	return i
}

func (i *RateLimiter) GetLimiter(identifier string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	v, exists := i.visitors[identifier]
	if !exists {
		r := rate.Limit(float64(RateLimitRequests) / RateLimitWindow.Seconds())
		limiter := rate.NewLimiter(r, RateLimitRequests)
		i.visitors[identifier] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// --- USAGE LOGS ---
type LogEntry struct {
	Timestamp  int64  `json:"timestamp"`
	Identifier string `json:"identifier"`
	Allowed    bool   `json:"allowed"`
}

var (
	usageLogs   []LogEntry
	usageLogsMu sync.RWMutex
)

const MaxLogs = 200

func addLog(identifier string, allowed bool) {
	usageLogsMu.Lock()
	defer usageLogsMu.Unlock()

	entry := LogEntry{
		Timestamp:  time.Now().UnixMilli(),
		Identifier: identifier,
		Allowed:    allowed,
	}

	// Prepend so the newest logs are at the beginning
	usageLogs = append([]LogEntry{entry}, usageLogs...)
	if len(usageLogs) > MaxLogs {
		usageLogs = usageLogs[:MaxLogs]
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/admin/logs", adminLogsHandler)

	// Initialize the rate limiter
	limiter := NewRateLimiter()

	// Rate Limiting Check Endpoint
	mux.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			Identifier string `json:"identifier"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Identifier == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request, missing 'identifier'"})
			return
		}

		userLimiter := limiter.GetLimiter(payload.Identifier)
		allowed := userLimiter.Allow()

		addLog(payload.Identifier, allowed)

		w.Header().Set("Content-Type", "application/json")
		if !allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"allowed": false,
				"error":   "Too Many Requests",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"allowed": true,
		})
	})

	// Serve static files for the Web Dashboard on the root path
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Rate Limiting Microservice is running on port %s", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func adminLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usageLogsMu.RLock()
	defer usageLogsMu.RUnlock()

	if usageLogs == nil {
		json.NewEncoder(w).Encode([]LogEntry{})
		return
	}
	json.NewEncoder(w).Encode(usageLogs)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "Rate Limiting Microservice is running",
	})
}

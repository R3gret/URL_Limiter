# Centralized Rate Limiting Microservice

A high-performance, language-agnostic Rate Limiting Microservice built in Go. This service provides a centralized architecture to manage and track API usage limits across multiple distributed systems using an incredibly fast, in-memory **Token Bucket** algorithm.

## 🌟 Key Features

- **Standalone Architecture:** Designed to be decoupled from your main application, allowing services written in Node.js, Python, PHP, or any other language to securely verify rate limits via a simple HTTP request.
- **Zero Dependencies:** Runs entirely on Go's standard library and the official `golang.org/x/time/rate` package. No external Redis or database setup required.
- **Dynamic Live Dashboard:** Features a premium, glassmorphism web dashboard built with HTML, CSS, and Vanilla JavaScript that displays real-time request logs without polling databases.
- **High Concurrency:** Utilizes Go routines and thread-safe Mutexes to handle thousands of concurrent rate-limit checks instantly.
- **Docker Ready:** Includes an optimized, multi-stage `Dockerfile` producing a tiny Alpine-based image ready for cloud deployment.

## 🚀 Getting Started

### Prerequisites
- [Go 1.21+](https://go.dev/doc/install) (if running directly)
- [Docker](https://docs.docker.com/get-docker/) (if running via container)

### Running Locally

1. Clone the repository and navigate into the directory.
2. Download dependencies:
   ```bash
   go mod download
   ```
3. Start the server:
   ```bash
   go run main.go
   ```
4. Open the interactive dashboard in your browser: `http://localhost:8080`

### Running with Docker

1. Build the production image:
   ```bash
   docker build -t rate-limiter-service .
   ```
2. Run the container:
   ```bash
   docker run -d -p 8080:8080 --name rate-limiter rate-limiter-service
   ```

## 🔌 API Documentation

### Endpoint: `POST /api/check`

Verifies if the given identifier has exceeded the rate limit.

**Request Body (JSON):**
```json
{
  "identifier": "user_123_or_ip_address"
}
```

**Response (200 OK):**
```json
{
  "allowed": true
}
```

**Response (429 Too Many Requests):**
```json
{
  "allowed": false,
  "error": "Too Many Requests"
}
```

## 💻 Integration Examples

Because this is a standalone REST microservice, integration takes only a few lines of code in your backend systems. 

**Node.js / Express Example:**
```javascript
const checkRateLimit = async (userId) => {
  const response = await fetch('http://localhost:8080/api/check', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ identifier: userId })
  });
  
  const data = await response.json();
  if (!data.allowed) {
    throw new Error("Rate limit exceeded!");
  }
};
```

**Python / FastAPI Example:**
```python
import requests

def check_rate_limit(user_id):
    response = requests.post(
        'http://localhost:8080/api/check',
        json={'identifier': user_id}
    )
    
    data = response.json()
    if not data.get('allowed'):
        raise Exception("Rate limit exceeded!")
```

## ⚙️ Configuration
By default, the server runs on port `8080`. You can override this by passing the `PORT` environment variable:
```bash
PORT=5000 go run main.go
```

The rate limit is currently set to **5 requests per 10 seconds**. To adjust this, modify the `RateLimitRequests` and `RateLimitWindow` constants inside `main.go`.

## 📜 License
This project is open-source and available under the MIT License.

# Step 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application (disabling CGO for a statically linked binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway main.go

# Step 2: Create a minimal production image
FROM alpine:latest

# Add certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/api-gateway .

# Copy the static files for the dashboard
COPY --from=builder /app/static ./static

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./api-gateway"]

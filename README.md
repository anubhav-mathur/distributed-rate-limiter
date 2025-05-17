# Distributed Rate Limiter (gRPC + Redis + Prometheus)

This is a distributed rate limiting service built in Go. It enforces per-user request quotas across multiple service instances using a Redis-backed token bucket algorithm. The system is observable through Prometheus metrics and a Grafana dashboard.

---

## Features

- gRPC API with endpoints for request validation and usage tracking
- Redis-based token bucket implementation (with Lua scripting for atomicity)
- Supports multi-instance deployment (shared quota enforcement)
- Prometheus metrics per instance (rate, status, latency)
- Grafana dashboard integration
- Configurable gRPC and metrics ports for distributed testing

---

## Tech Stack

- Go (gRPC server, Redis integration)
- Redis (shared coordination layer)
- gRPC + Protocol Buffers (RPC communication)
- Prometheus (metrics collection)
- Grafana (visualization)
- Docker Compose (for monitoring setup)

---

## Project Structure

distributed-rate-limiter/  
├── cmd/                  gRPC server entrypoint  
├── internal/  
│   ├── api/              gRPC handler logic  
│   ├── limiter/          Optional in-memory token bucket implementation  
│   ├── store/            Redis-based rate limiter  
│   └── metrics/          Prometheus metric definitions  
├── proto/                Protobuf definitions and generated files  
├── monitoring/           Prometheus and Grafana setup  
└── README.md

---

## API Endpoints

### AllowRequest

Validates if a request is allowed under the user's current rate limit.

**Request:**
    {
      "user_id": "user123",
      "path": "/login"
    }

**Response:**
    {
      "allowed": true,
      "reason": "Request allowed"
    }

---

### GetUsage

Returns the number of requests used and allowed for a given user.

**Request:**
    {
      "user_id": "user123"
    }

**Response:**
    {
      "requests_used": 3,
      "requests_allowed": 5
    }

---

## Getting Started

### 1. Run Redis

    docker run -d --name rate-limiter-redis -p 6379:6379 redis

---

### 2. Start the gRPC Server

    PORT=50051 METRICS_PORT=2112 go run cmd/server/main.go

To simulate multiple distributed instances:

    PORT=50052 METRICS_PORT=2113 go run cmd/server/main.go

---

### 3. Test with Postman or grpcurl

Postman: Use the gRPC tab to call `RateLimiter/AllowRequest` and `RateLimiter/GetUsage`.

grpcurl example:

    grpcurl -plaintext -d '{"user_id":"user123", "path":"/test"}' localhost:50051 limiter.RateLimiter/AllowRequest

---

## Monitoring with Prometheus and Grafana

### 1. Start Services

    cd monitoring
    docker compose up

### 2. Access UIs

- Prometheus: http://localhost:9090  
- Grafana: http://localhost:3000 (login: admin / admin)

### 3. Example Grafana Queries

| Metric | Description |
|--------|-------------|
| rate_limiter_requests_total | Total requests per user/status |
| sum by(status) (rate(rate_limiter_requests_total[1m])) | Allowed vs denied |
| histogram_quantile(0.95, rate(redis_latency_seconds_bucket[1m])) | Redis latency (95th percentile) |
| sum by(user) (rate(rate_limiter_requests_total[1m])) | Requests per user |

---

## Implementation Notes

- Token bucket: 5-token capacity, refills 1 token every 2 seconds
- Redis stores tokens and last refill time
- Redis commands are atomic via Lua
- Stateless servers; safe for horizontal scaling

---

## Testing Distributed Coordination

1. Run two gRPC server instances (different ports)
2. Send requests to both with the same user ID
3. Observe consistent global rate limiting enforced by Redis

---


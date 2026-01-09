---
code: GODEM-005
status: Proposed
dateCreated: 2026-01-09T19:59:12.332Z
type: Feature Enhancement
priority: Medium
dependsOn: GODEM-004
---

# Update health endpoint to check Redis and MongoDB connections

## 1. Description

### Problem Statement
The current `/health` endpoint only returns a static `{"status": "ok"}`. It doesn't verify that dependent services (MongoDB, Redis) are actually reachable.

### Current State
```json
{"status": "ok"}
```

### Desired State
```json
{
  "status": "ok",
  "services": {
    "mongodb": "connected",
    "redis": "connected"
  }
}
```

If a service is down:
```json
{
  "status": "degraded",
  "services": {
    "mongodb": "connected",
    "redis": "error: connection refused"
  }
}
```

### Business/Technical Justification
- Health checks are critical for container orchestration (Kubernetes, Docker Swarm)
- Enables monitoring systems to detect service degradation
- Supports graceful handling of partial outages

## 2. Rationale

### Why This Change Is Necessary
- Current health check gives false positives when dependencies are down
- Operators need visibility into service connectivity
- Required for proper load balancer health probes

### What It Accomplishes
- Accurate health status reflecting actual service state
- Quick identification of connectivity issues
- Foundation for future monitoring integration

## 3. Solution Analysis

### Selected Approach

| Aspect | Decision |
|--------|----------|
| Response format | JSON with nested services object |
| Status values | `ok`, `degraded`, `unhealthy` |
| HTTP status | 200 for ok/degraded, 503 for unhealthy |
| Timeout | 2s per service check |
| Parallel checks | Yes, check MongoDB and Redis concurrently |

### Status Logic

| MongoDB | Redis | Status | HTTP Code |
|---------|-------|--------|----------|
| connected | connected | `ok` | 200 |
| connected | error | `degraded` | 200 |
| error | connected | `degraded` | 200 |
| error | error | `unhealthy` | 503 |

## 4. Implementation Specification

### Files to Modify

| File | Changes |
|------|--------|
| `internal/handler/handler.go` | Add MongoDB/Redis clients as dependencies |
| `internal/handler/health.go` | Implement service checks |
| `features/health.feature` | Add scenarios for degraded/unhealthy states |

### Updated Handler Struct

```go
type Handler struct {
    mongo *mongodb.Client
    redis *redis.Client
}

func New(mongo *mongodb.Client, redis *redis.Client) *Handler
```

### Health Response Struct

```go
type healthResponse struct {
    Status   string            `json:"status"`
    Services map[string]string `json:"services"`
}
```

### Health Check Logic

```go
func (h *Handler) Health() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()

        services := make(map[string]string)
        
        // Check MongoDB (concurrent)
        // Check Redis (concurrent)
        
        status := determineOverallStatus(services)
        httpCode := statusToHTTPCode(status)
        
        // Return response
    }
}
```

## 5. Acceptance Criteria

- [ ] `/health` returns services status for MongoDB and Redis
- [ ] Status is `ok` when all services connected
- [ ] Status is `degraded` when one service is down
- [ ] Status is `unhealthy` (503) when all services are down
- [ ] Each service check has 2s timeout
- [ ] Checks run in parallel (not sequential)
- [ ] BDD tests cover all status scenarios
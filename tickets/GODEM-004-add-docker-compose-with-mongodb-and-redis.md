---
code: GODEM-004
status: Implemented
dateCreated: 2026-01-09T19:57:41.292Z
type: Feature Enhancement
priority: High
dependsOn: GODEM-001
blocks: GODEM-005
---

# Add Docker Compose with MongoDB and Redis

## 1. Description

### Problem Statement
The application currently runs as a single Docker container. To support real-world functionality, we need MongoDB for primary data storage and Redis for caching.

### Current State
- Single Dockerfile with Go app + Air hot-reload
- No database or cache services
- `make docker-run` runs single container

### Desired State
- Docker Compose orchestrating multiple containers (app, mongo, redis)
- Go app connects to MongoDB and Redis
- Persistent volumes for data durability
- Development workflow preserved (hot-reload still works)

### Business/Technical Justification
- MongoDB provides flexible document storage for application data
- Redis enables fast caching to improve performance
- Docker Compose simplifies multi-container development

## 2. Rationale

### Why This Change Is Necessary
- Real applications need persistent data storage
- Caching is essential for performance at scale
- Multi-service architecture is industry standard

### What It Accomplishes
- Enables data persistence across restarts
- Provides caching infrastructure
- Prepares architecture for production deployment

## 3. Solution Analysis

### Selected Approach

| Component | Choice | Rationale |
|-----------|--------|----------|
| Orchestration | Docker Compose | Standard for local multi-container dev |
| MongoDB | mongo:7 | Latest stable, good Go driver support |
| Redis | redis:7-alpine | Lightweight, latest stable |
| Go MongoDB driver | go.mongodb.org/mongo-driver | Official driver |
| Go Redis client | github.com/redis/go-redis/v9 | Official, well-maintained |
| Volumes | Named volumes | Persist data, easy cleanup |

### Alternatives Considered

| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Separate docker run commands | Simple | Manual, no networking | Rejected |
| Podman Compose | Rootless | Less common | Rejected |
| Embedded DB (SQLite) | No container needed | Not MongoDB | Rejected |

## 4. Implementation Specification

### New Files

```
go-demo/
├── docker-compose.yml           # Multi-container orchestration
├── internal/
│   ├── mongodb/
│   │   └── mongodb.go           # MongoDB client wrapper
│   └── redis/
│       └── redis.go             # Redis client wrapper
```

### docker-compose.yml Structure

```yaml
services:
  app:
    build: .
    ports:
      - "8080:80"
    volumes:
      - .:/app
    depends_on:
      - mongo
      - redis
    environment:
      - MONGO_URI=mongodb://mongo:27017
      - REDIS_ADDR=redis:6379

  mongo:
    image: mongo:7
    volumes:
      - mongo_data:/data/db
    ports:
      - "27017:27017"

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"

volumes:
  mongo_data:
  redis_data:
```

### Go Client Integration

#### MongoDB Client (`internal/mongodb/mongodb.go`)
```go
type Client struct {
    client *mongo.Client
    db     *mongo.Database
}

func New(uri, dbName string) (*Client, error)
func (c *Client) Close(ctx context.Context) error
func (c *Client) Ping(ctx context.Context) error
```

#### Redis Client (`internal/redis/redis.go`)
```go
type Client struct {
    client *redis.Client
}

func New(addr string) (*Client, error)
func (c *Client) Close() error
func (c *Client) Ping(ctx context.Context) error
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection string |
| `MONGO_DB` | `godemo` | Database name |
| `REDIS_ADDR` | `localhost:6379` | Redis address |

### Updated Makefile Targets

| Target | Command | Description |
|--------|---------|-------------|
| `make up` | `docker-compose up` | Start all services |
| `make down` | `docker-compose down` | Stop all services |
| `make logs` | `docker-compose logs -f` | Follow logs |
| `make clean-volumes` | `docker-compose down -v` | Remove data volumes |

### Health Endpoint Enhancement

Update `/health` to check MongoDB and Redis connectivity:
```json
{
  "status": "ok",
  "services": {
    "mongodb": "connected",
    "redis": "connected"
  }
}
```

## 5. Acceptance Criteria

- [ ] `docker-compose up` starts app, mongo, and redis containers
- [ ] Go app connects to MongoDB successfully
- [ ] Go app connects to Redis successfully
- [ ] `/health` reports status of all services
- [ ] Data persists after `docker-compose down` and `up`
- [ ] Hot-reload still works for Go code changes
- [ ] BDD tests pass with new infrastructure
- [ ] `make up` / `make down` commands work
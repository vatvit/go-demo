---
code: GODEM-001
status: Implemented
dateCreated: 2026-01-09T17:07:09.343Z
type: Architecture
priority: High
---

# Create Go boilerplate project with Docker and HTTP server
## 1. Description

### Problem Statement
The project needs a foundational Go application structure that supports local development with Docker. Currently, there is no codebase - only MDT ticket configuration exists.

### Current State
- Empty project with only `.mdt-config.toml` and `tickets/` directory
- No Go code, no Docker configuration

### Desired State
- Working Go HTTP server boilerplate
- Docker-based local development environment with hot-reload
- Simple, maintainable project structure

### Business/Technical Justification
This foundational setup enables rapid development of future features with consistent tooling and containerized environment.

## 2. Rationale

### Why This Change Is Necessary
- Establishes consistent development environment across machines
- Docker ensures reproducible builds and deployments
- Hot-reload improves developer experience during local development

### What It Accomplishes
- Provides a working starting point for Go development
- Standardizes the development workflow
- Enables future containerized deployments

### Alignment with Project Goals
- Focus on local development environment as primary use case
- Keep setup minimal and maintainable

## 3. Solution Analysis

### Selected Approach
- **Go Version:** 1.23 (latest stable)
- **HTTP Framework:** Standard library `net/http` (no external dependencies)
- **Docker Setup:** Single Dockerfile (no docker-compose)
- **Hot Reload:** Air (cosmtrek/air) for auto-rebuild on file changes
- **Port:** 80 inside container (mapped to host port via docker run)
- **Build Tooling:** Makefile for common commands

### Alternatives Considered
| Option | Pros | Cons | Decision |
|--------|------|------|----------|
| Gin/Chi/Echo framework | More features, routing | External dependency | Rejected - keep it simple |
| Docker Compose | Multi-service orchestration | Overkill for single service | Rejected |
| No hot-reload | Simpler setup | Poor DX | Rejected |

## 4. Implementation Specification
### Project Structure
```
go-demo/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   └── handler/
│       └── health.go        # Health check handler
├── .air.toml                 # Air hot-reload configuration
├── Dockerfile                # Container definition
├── Makefile                  # Build commands
├── go.mod                    # Go module definition
└── CLAUDE.md                 # Project instructions
```

### Technical Details

1. **main.go** - HTTP server setup
   - Listen on port 80
   - Register `/health` endpoint
   - Graceful shutdown handling

2. **Dockerfile** - Multi-stage build
   - Build stage: Go 1.23 alpine
   - Dev stage: Include Air for hot-reload
   - Mount source code as volume for development

3. **Makefile** targets:
   - `make build` - Build Go binary
   - `make run` - Run locally without Docker
   - `make docker-build` - Build Docker image
   - `make docker-run` - Run in Docker with hot-reload
   - `make test` - Run tests

4. **.air.toml** - Hot-reload config
   - Watch `.go` files
   - Rebuild on changes
   - Exclude `tmp/`, `vendor/`

### Endpoints
| Method | Path | Response | Purpose |
|--------|------|----------|--------|
| GET | `/health` | `{"status": "ok"}` | Liveness probe |

### HTTP Server Requirements

| ID | Type | Requirement |
|----|------|-------------|
| REQ-001 | Ubiquitous | The server shall listen on port 80 inside the container |
| REQ-002 | Ubiquitous | The server shall be accessible from host via port 8080 |
| REQ-003 | Event-driven | When a GET request is made to `/health`, the server shall return HTTP 200 with JSON body `{"status": "ok"}` |
| REQ-004 | Event-driven | When a request is made to an undefined route, the server shall return HTTP 404 with JSON body `{"error": "not found"}` |
| REQ-005 | Event-driven | When SIGTERM or SIGINT signal is received, the server shall initiate graceful shutdown |
| REQ-006 | State-driven | While shutting down, the server shall wait up to 30 seconds for in-flight requests to complete |
| REQ-007 | Unwanted | If shutdown timeout (30s) is exceeded, the server shall force terminate |

### Docker Development Environment Requirements

| ID | Type | Requirement |
|----|------|-------------|
| REQ-008 | Ubiquitous | The Dockerfile shall use Go 1.23 as the base image |
| REQ-009 | Ubiquitous | The Docker image shall include Air for hot-reload capability |
| REQ-010 | Event-driven | When a `.go` file is modified, Air shall trigger automatic rebuild and restart |
| REQ-011 | Ubiquitous | The source code shall be mounted as a volume for development |

### Build Tooling Requirements

| ID | Type | Requirement |
|----|------|-------------|
| REQ-012 | Ubiquitous | The Makefile shall provide `build` target to compile Go binary |
| REQ-013 | Ubiquitous | The Makefile shall provide `run` target to run locally without Docker |
| REQ-014 | Ubiquitous | The Makefile shall provide `docker-build` target to build Docker image |
| REQ-015 | Ubiquitous | The Makefile shall provide `docker-run` target to run with hot-reload |
| REQ-016 | Ubiquitous | The Makefile shall provide `test` target to run tests |
| REQ-017 | Ubiquitous | The Makefile shall provide `clean` target to remove binaries, tmp/, and Docker images |

### Constraints

| ID | Constraint |
|----|------------|
| CON-001 | No external Go dependencies for HTTP handling (stdlib only) |
| CON-002 | No request logging |
| CON-003 | Standard Go project layout (cmd/, internal/) |
## 5. Acceptance Criteria
- [ ] `go build` completes without errors
- [ ] `make docker-build` creates working Docker image
- [ ] `make docker-run` starts container with hot-reload
- [ ] `curl localhost:<port>/health` returns `{"status": "ok"}`
- [ ] Modifying Go code triggers automatic rebuild in Docker
- [ ] Project follows standard Go project layout
- [ ] No external Go dependencies (stdlib only for HTTP)

### Test Framework
- **Framework:** Godog (Cucumber for Go)
- **Location:** `features/` directory
- **Initial State:** RED (tests written before implementation)

### Feature: Health Check Endpoint

```gherkin

## 6. Architecture

### Project Structure (Detailed)

```
go-demo/
├── cmd/
│   └── server/
│       └── main.go              # Entry point: load config, create server, run
├── internal/
│   ├── config/
│   │   └── config.go            # Config struct, load from environment
│   ├── handler/
│   │   ├── handler.go           # Handler struct with dependencies
│   │   ├── health.go            # GET /health
│   │   └── notfound.go          # 404 JSON handler
│   └── server/
│       └── server.go            # Server struct, setup routes, graceful shutdown
├── features/                     # BDD tests (created after boilerplate)
│   ├── health.feature
│   ├── errors.feature
│   ├── shutdown.feature
│   ├── docker.feature
│   └── steps/
│       └── *.go
├── .air.toml                     # Air hot-reload configuration
├── Dockerfile                    # Multi-stage: build + dev with Air
├── Makefile                      # Build commands
├── go.mod                        # Go module definition
├── go.sum                        # Dependency checksums
└── CLAUDE.md                     # Project instructions
```

### Component Responsibilities

#### `cmd/server/main.go`
- Load configuration from environment
- Create Server instance
- Start HTTP server
- Handle OS signals for shutdown

#### `internal/config/config.go`
```go
type Config struct {
    Port            string        // HTTP port (default: "80")
    ShutdownTimeout time.Duration // Graceful shutdown timeout (default: 30s)
}

func Load() (*Config, error)  // Load from environment variables
```

**Environment Variables:**
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `80` | HTTP server port |
| `SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |

#### `internal/server/server.go`
```go
type Server struct {
    config  *config.Config
    handler *handler.Handler
    server  *http.Server
}

func New(cfg *config.Config) *Server           // Constructor
func (s *Server) Run(ctx context.Context) error // Start and block
func (s *Server) Shutdown(ctx context.Context) error // Graceful stop
```

#### `internal/handler/handler.go`
```go
type Handler struct {
    // Future dependencies go here (db, logger, etc.)
}

func New() *Handler
func (h *Handler) Routes() http.Handler  // Returns configured mux
```

#### `internal/handler/health.go`
```go
func (h *Handler) Health() http.HandlerFunc
// Returns: 200 {"status": "ok"}
```

#### `internal/handler/notfound.go`
```go
func (h *Handler) NotFound() http.HandlerFunc
// Returns: 404 {"error": "not found"}
```

### Data Flow

```
Request → http.Server → ServeMux → Handler method → JSON Response
                              ↓
                        NotFound (if no match)
```

### Shutdown Flow

```
SIGTERM/SIGINT → main.go catches signal
              → ctx cancel
              → server.Shutdown(ctx with 30s timeout)
              → wait for in-flight requests
              → exit 0
```

### Docker Architecture

```
┌─────────────────────────────────────────┐
│ Host Machine                            │
│  ┌───────────────────────────────────┐  │
│  │ Docker Container                  │  │
│  │  ┌─────────────┐                  │  │
│  │  │ Air         │ watches .go      │  │
│  │  │  ↓ rebuild  │                  │  │
│  │  │ Go Binary   │ :80              │  │
│  │  └─────────────┘                  │  │
│  │       ↑                           │  │
│  │  /app (volume mount)              │  │
│  └───────────────────────────────────┘  │
│       ↑ :8080 → :80                     │
│  localhost:8080                         │
└─────────────────────────────────────────┘
```

## 7. Tasks

### Implementation Order

| # | Task | Files | Depends On |
|---|------|-------|------------|
| 1 | Initialize Go module | `go.mod` | - |
| 2 | Create config package | `internal/config/config.go` | 1 |
| 3 | Create handler package | `internal/handler/handler.go`, `health.go`, `notfound.go` | 1 |
| 4 | Create server package | `internal/server/server.go` | 2, 3 |
| 5 | Create main entry point | `cmd/server/main.go` | 4 |
| 6 | Create Air config | `.air.toml` | - |
| 7 | Create Dockerfile | `Dockerfile` | 6 |
| 8 | Create Makefile | `Makefile` | 7 |
| 9 | Update CLAUDE.md | `CLAUDE.md` | 8 |
| 10 | Verify: manual test | - | 9 |

### Task Details

#### Task 1: Initialize Go module
```bash
go mod init github.com/vatvit/go-demo
```

#### Task 2: Create config package
- Create `internal/config/config.go`
- Define `Config` struct with `Port` and `ShutdownTimeout`
- Implement `Load()` function reading from env vars
- Defaults: PORT=80, SHUTDOWN_TIMEOUT=30s

#### Task 3: Create handler package
- `internal/handler/handler.go` - Handler struct + New() + Routes()
- `internal/handler/health.go` - Health() returns 200 `{"status": "ok"}`
- `internal/handler/notfound.go` - NotFound() returns 404 `{"error": "not found"}`

#### Task 4: Create server package
- `internal/server/server.go`
- Server struct with config, handler, http.Server
- New() constructor
- Run(ctx) - start server, block until ctx done
- Shutdown(ctx) - graceful shutdown with timeout

#### Task 5: Create main entry point
- `cmd/server/main.go`
- Load config
- Create server
- Setup signal handling (SIGTERM, SIGINT)
- Run server, shutdown on signal

#### Task 6: Create Air config
- `.air.toml`
- Watch `cmd/`, `internal/` for .go files
- Exclude `tmp/`, `vendor/`, `features/`
- Build command: `go build -o tmp/main cmd/server/main.go`
- Run command: `./tmp/main`

#### Task 7: Create Dockerfile
- Base: `golang:1.23-alpine`
- Install Air: `go install github.com/air-verse/air@latest`
- WORKDIR /app
- EXPOSE 80
- CMD for dev: `air -c .air.toml`

#### Task 8: Create Makefile
Targets:
- `build` - go build
- `run` - go run
- `test` - go test
- `docker-build` - docker build
- `docker-run` - docker run with volume mount, port 8080:80
- `clean` - remove tmp/, binaries, docker image
- `bdd` - placeholder for godog tests

#### Task 9: Update CLAUDE.md
- Add build/run instructions
- Document Makefile targets
- Add development workflow

#### Task 10: Verify
- `make docker-build`
- `make docker-run`
- `curl localhost:8080/health` → `{"status": "ok"}`
- `curl localhost:8080/foo` → `{"error": "not found"}`
- Modify .go file → verify hot-reload

## 8. Reflection

### What Went Well

1. **MDT Workflow** - Following the structured workflow (ticket → requirements → BDD → architecture → tasks → implement) provided clear documentation and traceability.

2. **Docker-only development** - No local Go installation required. All commands run via Docker, ensuring consistent environment.

3. **BDD-first approach** - Writing feature files before implementation clarified expected behavior, even though this was a boilerplate project.

4. **Incremental tech debt resolution** - Identifying and fixing debt immediately after implementation prevented accumulation.

### Challenges Encountered

1. **Air version incompatibility** - Latest Air (v1.63.6) requires Go 1.25, but we use Go 1.23. Fixed by pinning to v1.61.7.

2. **Chicken-egg with BDD** - Can't run BDD tests without Go module and dependencies. For boilerplate projects, some workflow steps need reordering.

3. **Custom 404 handler pattern** - Go's `http.ServeMux` doesn't have a native "catch-all" for 404. Used `mux.Handler(r)` pattern check which works but is slightly unconventional.

### Lessons Learned

1. **Pin dependency versions** - Always pin versions in Dockerfile to avoid breaking changes (Air, Go).

2. **Boilerplate is an exception** - For greenfield projects, the strict BDD-before-implementation flow needs adaptation since there's no code to test against initially.

3. **Docker-first simplifies onboarding** - New developers only need Docker, not Go toolchain.

4. **Makefile as interface** - All commands via Makefile creates consistent developer experience regardless of underlying tools.

### Technical Decisions Worth Remembering

| Decision | Rationale |
|----------|-----------|
| Go 1.23 + Air v1.61.7 | Compatibility constraint |
| `net/http` stdlib only | Simplicity, no external deps |
| Port 80 internal, 8080 external | Standard container pattern |
| Config from env vars | 12-factor app compliance |
| Server struct with deps | Testability, extensibility |
| Cached 404 handler | Avoid per-request allocation |

### Metrics

- **Files created:** 14
- **Lines of Go code:** ~150
- **BDD scenarios:** 3 implemented, 4 pending (GODEM-003)
- **Tech debt items:** 4 fixed, 1 moved to backlog
# features/health.feature
Feature: Health Check Endpoint
  As a system operator
  I want a health check endpoint
  So that I can monitor if the server is running

  Scenario: Server returns healthy status
    Given the server is running
    When I send a GET request to "/health"
    Then the response status code should be 200
    And the response content type should be "application/json"
    And the response body should be JSON:
      """
      {"status": "ok"}
      """
```

### Feature: Error Handling

```gherkin
# features/errors.feature
Feature: Error Handling
  As an API consumer
  I want consistent error responses
  So that I can handle errors programmatically

  Scenario: Accessing undefined route returns 404
    Given the server is running
    When I send a GET request to "/undefined-route"
    Then the response status code should be 404
    And the response content type should be "application/json"
    And the response body should be JSON:
      """
      {"error": "not found"}
      """

  Scenario: Accessing undefined route with POST returns 404
    Given the server is running
    When I send a POST request to "/undefined-route"
    Then the response status code should be 404
    And the response content type should be "application/json"
```

### Feature: Graceful Shutdown

```gherkin
# features/shutdown.feature
Feature: Graceful Shutdown
  As a system operator
  I want the server to shutdown gracefully
  So that in-flight requests are not dropped

  Scenario: Server handles SIGTERM gracefully
    Given the server is running
    And a slow request is in progress
    When SIGTERM signal is sent to the server
    Then the in-flight request should complete successfully
    And the server should stop accepting new connections
    And the server should exit within 30 seconds

  Scenario: Server handles SIGINT gracefully
    Given the server is running
    When SIGINT signal is sent to the server
    Then the server should exit within 30 seconds
    And the exit code should be 0
```

### Feature: Docker Development Environment

```gherkin
# features/docker.feature
Feature: Docker Development Environment
  As a developer
  I want to run the server in Docker with hot-reload
  So that I can develop efficiently

  Scenario: Docker container starts successfully
    Given the Docker image is built
    When I run the container with port 8080 mapped to 80
    Then the container should be running
    And I should be able to access "http://localhost:8080/health"

  Scenario: Hot-reload triggers on file change
    Given the Docker container is running with source mounted
    When I modify a .go file in the source directory
    Then Air should detect the change
    And the server should rebuild automatically
    And the server should restart within 10 seconds
```

### Step Definitions Structure

```
features/
├── health.feature
├── errors.feature
├── shutdown.feature
├── docker.feature
└── steps/
    ├── server_steps.go      # Given the server is running, etc.
    ├── http_steps.go        # When I send a GET request, etc.
    ├── response_steps.go    # Then the response status code, etc.
    ├── signal_steps.go      # When SIGTERM signal is sent, etc.
    └── docker_steps.go      # Given the Docker image is built, etc.
```

### Makefile BDD Targets

| Target | Description |
|--------|-------------|
| `make bdd` | Run all BDD tests |
| `make bdd-health` | Run health feature only |
| `make bdd-docker` | Run docker feature only |
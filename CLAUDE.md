# Go Demo Project

## Project Overview

This is a Go demonstration project with Docker-based development environment.

## Quick Start

```bash
# Build and run with hot-reload
make docker-run

# Access the server
curl http://localhost:8080/health
```

## Ticket Management

This project uses **MDT (Markdown Documentation Tickets)** for work management.

- **Project Code:** GODEM
- **Tickets Location:** `tickets/`
- **CR Naming:** GODEM-XXX

### MDT Workflow

Follow this workflow for new features (see `tickets/GODEM-001` as reference):

1. `/mdt:ticket-creation` - Create CR with full specification
2. `/mdt:requirements` - Document EARS format requirements
3. `/mdt:bdd` - Write feature files (Gherkin) - tests should be RED initially
4. `/mdt:architecture` - Define structure, components, data flow
5. `/mdt:tasks` - Create implementation task list
6. `/mdt:implement` - Execute tasks (TDD: RED → GREEN)
7. `/mdt:tech-debt` - Analyze and document any debt
8. `/mdt:reflection` - Document learnings

**Note:** For boilerplate/greenfield projects, BDD tests may need to be implemented alongside the code rather than strictly before.

### Working with Tickets

1. Before starting work, check for existing CRs using MDT tools
2. Create new CRs for any significant changes (features, bugs, technical debt)
3. Update CR status as work progresses:
   - `Proposed` → `Approved` → `In Progress` → `Implemented`
4. Link related CRs using `dependsOn`, `blocks`, and `relatedTickets` attributes

### CR Types

- **Architecture** - System design changes
- **Feature Enhancement** - New functionality
- **Bug Fix** - Defect resolution
- **Technical Debt** - Code quality improvements
- **Documentation** - Project documentation

## Development

### Prerequisites

- Docker (no local Go installation required)

### Important Constraints

- **No local commands:** Do not run `go`, `npm`, or other tools directly on host. Use Docker for everything.
- **Air version:** Pinned to v1.61.7 (latest requires Go 1.25+, we use Go 1.23)
- **Stdlib only:** HTTP handling uses `net/http` only, no external frameworks

### Makefile Targets

| Command | Description |
|---------|-------------|
| `make help` | Show all available targets |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run with hot-reload (port 8080) |
| `make docker-stop` | Stop running container |
| `make build` | Build Go binary (via Docker) |
| `make run` | Run without hot-reload |
| `make test` | Run Go tests |
| `make bdd` | Run BDD/Godog tests |
| `make clean` | Remove artifacts and Docker image |

### Project Structure

```
go-demo/
├── cmd/server/main.go       # Entry point
├── internal/
│   ├── config/              # Configuration from env vars
│   ├── handler/             # HTTP handlers (health, notfound)
│   └── server/              # Server struct, lifecycle, graceful shutdown
├── features/                # BDD tests (Godog/Gherkin)
│   ├── *.feature            # Feature files
│   └── steps/               # Step definitions
├── tickets/                 # MDT tickets (GODEM-XXX)
├── .air.toml                # Hot-reload configuration
├── .gitignore               # Excludes tmp/, build artifacts
├── .dockerignore            # Excludes .git/, tickets/, *.md
├── Dockerfile               # Development container (Go 1.23 + Air)
├── Makefile                 # Build commands
└── go.mod                   # Go module (github.com/vatvit/go-demo)
```

### Architecture Patterns

| Pattern | Implementation |
|---------|----------------|
| Config from env | `internal/config/config.go` - `Load()` reads PORT, SHUTDOWN_TIMEOUT |
| Server struct | `internal/server/server.go` - Holds config, handler, http.Server |
| Handler struct | `internal/handler/handler.go` - Extensible for future deps (db, logger) |
| Cached handlers | NotFound handler cached in `Routes()` to avoid per-request allocation |
| Graceful shutdown | 30s timeout, handles SIGTERM/SIGINT |
| JSON error handling | All `json.Encode()` calls have error handling with fallback |

### Endpoints

| Method | Path | Response |
|--------|------|----------|
| GET | `/health` | `{"status": "ok"}` |
| * | `/*` | `{"error": "not found"}` (404) |

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `80` | HTTP server port |
| `SHUTDOWN_TIMEOUT` | `30s` | Graceful shutdown timeout |

## Development Guidelines

### Go Conventions

- Follow standard Go project layout (`cmd/`, `internal/`)
- Use `go fmt` for formatting
- Run `go vet` and `go test` before committing
- All commands run via Docker (no local Go required)
- Handle all errors - don't ignore return values

### Adding New Endpoints

1. Create handler method in `internal/handler/`
2. Register route in `handler.go` `Routes()` method
3. Add BDD scenario in `features/`
4. Run `make bdd` to verify

### Git Workflow

- Do NOT commit or push without explicit user approval
- Create descriptive commit messages referencing CR numbers (e.g., "GODEM-001: Add feature X")

## Known Issues / Backlog

| Ticket | Priority | Description |
|--------|----------|-------------|
| GODEM-003 | Low | Missing BDD scenarios for shutdown and docker hot-reload |

## Lessons Learned (from GODEM-001)

1. **Pin dependency versions** in Dockerfile to avoid breaking changes
2. **Docker-first** simplifies onboarding - only Docker required
3. **Makefile as interface** - consistent developer experience
4. **BDD for boilerplate** needs adaptation - can't test before code exists
5. **Tech debt review** immediately after implementation prevents accumulation

## Generated Code

- Do NOT modify generated code directly
- Modify templates and regenerate instead

---
code: GODEM-003
status: Proposed
dateCreated: 2026-01-09T19:27:13.301Z
type: Feature Enhancement
priority: Low
relatedTickets: GODEM-001,GODEM-002
---

# Implement missing BDD scenarios (shutdown, docker)

## 1. Description

GODEM-001 specification defined BDD scenarios for graceful shutdown and Docker hot-reload testing, but only health and error handling scenarios were implemented.

### Missing Scenarios

**shutdown.feature:**
- Server handles SIGTERM gracefully
- Server handles SIGINT gracefully
- In-flight requests complete before shutdown

**docker.feature:**
- Docker container starts successfully
- Hot-reload triggers on file change

## 2. Rationale

These scenarios validate critical functionality:
- Graceful shutdown prevents dropped requests in production
- Hot-reload is essential for developer experience

## 3. Implementation Specification

### Files to Create

```
features/
├── shutdown.feature      # New
├── docker.feature        # New
└── steps/
    ├── signal_steps.go   # New - SIGTERM/SIGINT handling
    └── docker_steps.go   # New - container lifecycle
```

### Technical Challenges

- **Shutdown tests:** Require spawning actual server process (not httptest)
- **Docker tests:** Require Docker-in-Docker or external test runner

## 4. Acceptance Criteria

- [ ] `features/shutdown.feature` implemented with passing tests
- [ ] `features/docker.feature` implemented with passing tests
- [ ] All 7 BDD scenarios from GODEM-001 spec pass
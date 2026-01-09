---
code: GODEM-002
status: Implemented
dateCreated: 2026-01-09T19:17:34.778Z
type: Technical Debt
priority: Low
relatedTickets: GODEM-001
---

# Tech debt from GODEM-001 boilerplate implementation
## 1. Description

### Problem Statement
During implementation of GODEM-001 (Go boilerplate), several minor technical debt items were identified that don't affect functionality but could be improved for code quality and maintainability.

### Current State
Boilerplate is functional but has minor code quality issues.

### Desired State
Clean codebase following Go best practices.

## 2. Rationale

These items are low priority since they don't affect functionality, but addressing them improves:
- Build performance (dockerignore)
- Repository cleanliness (gitignore)
- Error resilience (JSON encoding)
- Test coverage (missing BDD scenarios)

## 3. Tech Debt Items
### TD-1: Missing .gitignore ✅ FIXED
**Location:** Project root
**Issue:** No `.gitignore` file. `tmp/` directory with compiled binary is not excluded.
**Impact:** Low - may accidentally commit build artifacts
**Fix:** Add `.gitignore` with `tmp/`, `*.exe`, `.DS_Store`, etc.

### TD-2: Missing .dockerignore ✅ FIXED
**Location:** Project root
**Issue:** No `.dockerignore` file. Docker build copies unnecessary files (`.git/`, `tmp/`, `tickets/`).
**Impact:** Low - slower builds, larger build context
**Fix:** Add `.dockerignore` excluding `.git`, `tmp`, `tickets`, `*.md`

### TD-3: JSON encoding error not handled ✅ FIXED
**Location:** `internal/handler/health.go:16`, `internal/handler/notfound.go:16`
**Issue:** `json.NewEncoder(w).Encode()` return value (error) is ignored.
**Impact:** Low - if encoding fails, response is incomplete but this is unlikely for simple structs
**Fix:** Handle error with fallback to plain text error response

### TD-4: NotFound handler recreated per request ✅ FIXED
**Location:** `internal/handler/handler.go:20`
**Issue:** `h.NotFound()()` creates a new `http.HandlerFunc` on every 404 request.
**Impact:** Very Low - minor allocation overhead
**Fix:** Cache the handler: `notFoundHandler := h.NotFound()` then use `notFoundHandler(w, r)`

### TD-5: Missing BDD scenarios (shutdown, docker) → GODEM-003
**Status:** Moved to separate ticket GODEM-003 (Feature Enhancement, Low priority)
## 4. Implementation Specification

### Priority Order
1. TD-1 + TD-2 (gitignore/dockerignore) - Quick wins
2. TD-3 + TD-4 (code quality) - Minor fixes
3. TD-5 + TD-6 (test coverage) - More effort

### Estimated Effort
- TD-1: 5 min
- TD-2: 5 min
- TD-3: 10 min
- TD-4: 5 min
- TD-5: 30 min
- TD-6: 45 min

## 5. Acceptance Criteria

- [ ] `.gitignore` excludes build artifacts
- [ ] `.dockerignore` reduces build context
- [ ] No unhandled errors in handler code
- [ ] NotFound handler is cached
- [ ] All BDD scenarios from spec are implemented
- [ ] Unit test coverage > 70%

## 6. Resolution Log
### Fixed Items

| ID | Status | Date | Notes |
|----|--------|------|-------|
| TD-1 | ✅ Fixed | 2026-01-09 | Added `.gitignore` |
| TD-2 | ✅ Fixed | 2026-01-09 | Added `.dockerignore` |
| TD-3 | ✅ Fixed | 2026-01-09 | Added error handling to JSON encode |
| TD-4 | ✅ Fixed | 2026-01-09 | Cached `notFoundHandler` in `Routes()` |
| TD-5 | → GODEM-003 | 2026-01-09 | Moved to separate low-priority ticket |
| TD-6 | ❌ Removed | 2026-01-09 | Not tech debt - unit tests are a feature |

### Summary
All tech debt items from GODEM-001 have been resolved.
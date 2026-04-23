# Modernize Codebase Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Update this Go API to current dependencies and idioms while improving structure, testability, and startup safety.

**Architecture:** Move the executable entrypoint to `cmd/activities`, split domain, storage, config, and HTTP handling into `internal/...`, and replace global initialization with explicit wiring. Use the standard library router and dependency injection so handlers can be tested without Firestore.

**Tech Stack:** Go 1.26, `net/http` `ServeMux`, Cloud Firestore client, table-driven tests with the standard testing package.

---

### Task 1: Lock the upgrade target and package layout

**Files:**
- Create: `cmd/activities/main.go`
- Create: `internal/config/config.go`
- Create: `internal/activity/types.go`
- Create: `internal/store/firestore.go`
- Create: `internal/httpapi/handler.go`
- Modify: `go.mod`

- [ ] Record the new package boundaries and raise the Go/toolchain target in `go.mod`.
- [ ] Remove obsolete routing dependencies if the standard library replaces them.

### Task 2: Write failing tests for validation and HTTP behavior

**Files:**
- Create: `internal/activity/validate_test.go`
- Create: `internal/httpapi/handler_test.go`

- [ ] Add validator tests for empty author/title/message cases.
- [ ] Add handler tests for `GET /health`, `GET /activities`, and one create/update flow against a fake store.
- [ ] Run targeted tests and confirm they fail before implementation.

### Task 3: Implement the refactor behind the tests

**Files:**
- Create: `internal/activity/validate.go`
- Create: `internal/activity/service.go`
- Create: `internal/httpapi/routes.go`
- Modify: `internal/store/firestore.go`

- [ ] Introduce typed config and explicit app startup.
- [ ] Replace package-global Firestore state with a store struct and context-aware methods.
- [ ] Replace `github.com/eelcoh/go-api` and `gorilla/mux` with `http.ServeMux` route patterns and `r.PathValue`.
- [ ] Keep existing endpoints, but allow internal API cleanup where it simplifies the code.

### Task 4: Verify and clean up repo metadata

**Files:**
- Modify: `Dockerfile`
- Modify: `Makefile`
- Modify: `AGENTS.md`

- [ ] Update build paths so the container still copies a binary named `activities`.
- [ ] Run `go test ./...`, `go test -race ./...`, and `go mod tidy`.
- [ ] Refresh contributor guidance if commands or config changed materially.

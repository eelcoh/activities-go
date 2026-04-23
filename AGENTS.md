# Repository Guidelines

## Project Structure & Module Organization
This repository is a small Go API service with a flat root-level layout. `main.go` bootstraps Firestore and starts the HTTP server, `routes.go` registers endpoints, `handlers.go` contains request handlers, `objects.go` manages Firestore reads and writes, and `types.go`/`validate.go` define domain models and validation rules. Build artifacts are produced as `./activities`. Container packaging lives in `Dockerfile`.

## Build, Test, and Development Commands
Use the standard Go toolchain unless you are publishing images.

- `go build .` or `make local`: build the local binary.
- `go run . 8080`: start the API on port `8080` for local checks.
- `go test ./...`: run all tests once test files exist.
- `make linux`: build a static Linux binary for container use.
- `make docker` or `make image`: build Docker images from the current binary.

## Coding Style & Naming Conventions
Follow default Go formatting: tabs for indentation, `gofmt` before committing, and grouped imports. Keep package names lowercase and unexported identifiers in `camelCase`. Use clear handler names such as `retrieveActivities` or `updateComment`; new tests should follow Go’s `TestXxx` naming. Keep new files in the repo root unless a real package boundary emerges.

## Testing Guidelines
There are currently no committed `*_test.go` files, so add tests alongside the code they cover. Prefer table-driven unit tests for validators and handler behavior. Run `go test ./...` before opening a PR. When adding Firestore-dependent logic, isolate pure functions where possible so tests do not require live cloud resources.

## Commit & Pull Request Guidelines
Recent commits use short, imperative subjects such as `added sorting order` and `simplified activities`. Keep commit messages concise, lower-case, and focused on one change. PRs should include a short summary, any required environment changes, and example requests or responses when API behavior changes.

## Security & Configuration Tips
The service expects `PROJECTID`, `SIGNINGKEY`, and `PASSPHRASE` in the environment before startup. Do not commit secrets or local credential files. Verify container and deployment changes against the scratch-based `Dockerfile`, which expects the compiled binary to be named `activities`.

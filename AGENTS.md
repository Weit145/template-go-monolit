# Repository Guidelines

## Project Structure & Module Organization

- `cmd/app`: application entry point and signal handling.
- `internal/app`: dependency wiring and resource lifecycle.
- `internal/domain`: business entities, validation, and domain errors.
- `internal/usecase`: business scenarios; `port` defines dependency contracts.
- `internal/repositories/storage/postgres`: PostgreSQL persistence and transaction support.
- `internal/transport/api`: HTTP handlers, DTOs, middleware, and WebSocket clients.
- `internal/events`: event distribution and WebSocket subscriptions.
- `migrations`: Goose SQL migrations.
- `docs/api/v1/swagger.yaml`: OpenAPI specification.

Place tests alongside the code they exercise.

## Build, Test, and Development Commands

Use Go compatible with `go.mod`, Task, and Docker Compose.

- `task build`: build images, start PostgreSQL, apply migrations, and start the service.
- `task up` / `task down`: start or stop the stack.
- `task app:build`: build the application into `bin/app`.
- `task test`: run `go test ./...`.
- `go vet ./...`: check for common Go mistakes.
- `task migrate:create -- create_products`: create a migration.
- `task migrate:up`: apply migrations using `PG_DSN`.

`task delete` also deletes the database volume.

## Coding Style & Naming Conventions

Format Go code with `gofmt`; use standard Go indentation and naming. Pass `context.Context` first. Preserve underlying errors with `%w` or `errors.Join`, and inspect them with `errors.Is` or `errors.As`.

Keep HTTP DTOs and SQL models separate from domain entities. Publish success events after transaction commit.

## Testing Guidelines

Use Go’s standard `testing` package. Name files `*_test.go` and functions `TestBehavior`; use table-driven tests for validation and error cases.

No tests or coverage threshold currently exist. Add focused regression tests for behavior changes. Verify transaction rollback and WebSocket lifecycle when changing those paths.

## Commit & Pull Request Guidelines

History uses short, action-oriented subjects such as `Fix context bug`; Conventional Commits are not required.

Keep commits focused. PRs should explain the change, reference relevant issues, and report verification commands. Update OpenAPI documentation when API behavior changes.

## Security & Configuration

Use `.env.example` as the configuration reference; never commit credentials or `.env`. Direct `go run` execution does not automatically load `.env`.

Create new migrations for deployed schema changes rather than rewriting applied migrations.

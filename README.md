# Link Shortener Development Guide

## Prerequisites
- Docker (Desktop 4.x or newer) with Docker Compose v2
- GNU Make 4.x+
- Go 1.22+ (only required for `go run` outside Docker)

## Environment Setup
1. Copy these variables into a `.env` file at the repository root (adjust as needed):
   ```env
   POSTGRES_DB=linksh
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_PORT=5432
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/linksh?sslmode=disable
   ```
2. Ensure the `.env` file is loaded by Docker Compose (Compose automatically reads it).
3. If you plan to run the Go app directly, export the same variables in your shell or use a tool like `direnv`.

## Starting Work
Most workflows happen through the Makefile wrappers (which call Docker Compose under the hood).

### Start the database only
- `make db-up`
  - Brings the Postgres container (`linksh-db`) up in detached mode.
  - Required before running migrations or the app.

### Apply database migrations
- `make migrate-up`
  - Runs `docker compose run --rm migrate up` against `DATABASE_URL`.
  - Use this after any schema change or the first time you bring the database up.
- `make migrate-down`
  - Rolls back the most recent migration if you need to revert a change.

### Start the application
- `make app-up`
  - Starts the API container (`linksh-app`) and streams logs in the foreground. Exit with `Ctrl+C` to stop just the app.
- `make dev`
  - Starts the database in detached mode, applies migrations, then runs the app with logs attached. Ideal for a clean daily start.
- `make dev-up`
  - Starts both `db` and `app` in detached mode. Follow with `make migrate-up` if needed.

If you prefer raw Docker Compose commands:
- `docker compose up -d db`
- `docker compose up --build app`

## Stopping and Restarting Services
- Stop only the app: `docker compose stop app` (or `Ctrl+C` if running attached).
- Stop only the database: `make db-down`.
- Restart the app (without rebuilding): `docker compose restart app`.
- Rebuild and restart the app: `make app-reset`.
- Restart just the database: `docker compose restart db`.
- Tear everything down (including networks, but keeping data volume): `make dev-down`.
- Fully reset the database (drops volume, reruns migrations): `make db-reset` (destructive!).

## Working with Migrations
- Apply the latest migrations: `make migrate-up` (database must be running).
- Roll back one step: `make migrate-down`.
- Create a new migration: `make migrate-create name=<migration_name>` (requires the `migrate` CLI via Docker as defined in the Makefile).
  - Example: `make migrate-create name=add_users_table`.
  - Edit the generated `.up.sql` and `.down.sql` files under `cmd/migrate/migrations/`, then run `make migrate-up` to apply.

## Running the Go app directly (optional)
If you want to bypass Docker for local development:
1. Ensure Postgres is running (either via Docker or locally) and `DATABASE_URL` points to it.
2. Run `go run ./cmd/main`.

## Original Notes
```
curl -X POST "localhost:8080/routemap?path=blsah&destination=google.com"
curl -X GET "localhost:8080/routemap?path=blsah"
```

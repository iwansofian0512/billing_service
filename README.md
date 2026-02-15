# Billing Engine API

A small billing engine written in Go. It manages borrowers, loans, billing schedules, and payments on a weekly schedule.

## Tech Stack

- Go with Gin HTTP framework
- PostgreSQL as the primary database
- Docker and Docker Compose for local development

## Getting Started

### Prerequisites

- Go toolchain installed
- Docker and Docker Compose

### Environment Variables

Use `.env.example` as a reference:

- `PORT` – HTTP port for the API server (default: `8080`)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` – Postgres connection settings used by the app and `config/db/postgres.go`.

Create a local `.env`:

```bash
cp .env.example .env
# adjust values as needed
```

## Makefile Commands

The project includes a `Makefile` with common tasks:

- `make build` – Build the Go binary into `bin/billing_service` (or `bin/$APP_NAME`).
- `make run` – Run the API server with `go run cmd/server/main.go`.
- `make test` – Run all Go tests (`go test ./...`).
- `make docker-build` – Build the Docker image for the application.
- `make docker-up` – Start the app and Postgres using `docker compose up`.
- `make docker-down` – Stop all Docker Compose services.
- `make cleanup-db` – Stop services and remove volumes (reset local database data).
- `make migrate-db` – Apply SQL migrations inside the Postgres container:
  - `migrations/init_schema.up.sql`
  - `migrations/sample_data.up.sql`

### Typical Local Flow

```bash
# build and run directly
make build
make run

# or use Docker
make docker-up
make migrate-db   # only needed the first time or after cleanup
```

## Project Structure

- `cmd/server` – application entrypoint, loads env, wiring, and graceful HTTP shutdown.
- `config/db` – PostgreSQL connection factory using `sqlx`.
- `internal/handler` – HTTP handlers and Gin router.
- `internal/service` – core business logic for borrowers, loans, and payments.
- `internal/repository` – data access layer for Postgres.
- `internal/model` – shared domain models and request/response payloads.
- `migrations` – SQL migrations for schema and sample data.
- `postman.json` – Postman collection with example API requests.

## Database Migrations

Migrations are plain SQL files:

- `migrations/init_schema.up.sql` – creates tables and enums.
- `migrations/sample_data.up.sql` – inserts example borrowers, loans, schedules, and payments.

The `make migrate-db` command uses `psql` inside the `db` service (from `docker-compose.yml`) to apply these files to the `billing_service` database.

## API Overview

Router: `internal/handler/router.go`

### Borrowers

- `POST /api/v1/borrowers` – create a borrower.
- `GET /api/v1/borrowers?borrower_id={id}&page={n}&page_size={m}` – list loans for a borrower with basic pagination.

### Loans

- `POST /api/v1/loans` – create a new loan for a borrower and generate weekly billing schedules.

### Payments

- `POST /api/v1/payment` – make a weekly payment against a loan.

The payment logic lives in `internal/service/payment_service/payment_service.go` and updates both the billing schedule and loan status.

## AI USAGE

AI usage for non functional code like README, sample_data.up.sql, Makefile, postman.json and fixing some unit test. Functional code written manualy with some refference from my previous work experience.

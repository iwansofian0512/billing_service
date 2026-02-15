ifneq (,$(wildcard .env))
include .env
export $(shell sed -n 's/^\([A-Za-z_][A-Za-z0-9_]*\)=.*/\1/p' .env)
endif

APP_NAME ?= billing_service

# Build the Go binary into bin/ with the configured app name
build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

# Run the API server directly with go run
run:
	go run cmd/server/main.go

# Run all Go tests in the module
test:
	go test ./...

# Build the Docker image for the application
docker-build:
	docker build -t $(APP_NAME) .

# Start application and database using docker compose
docker-up:
	docker compose up 

# Stop all docker compose services
docker-down:
	docker compose down

# Drop all application database objects using SQL down migration
cleanup-db:
	docker compose exec -T db psql -U ${DB_USER} -d ${DB_NAME} < migrations/init_schema.down.sql

# Apply schema and sample data migrations into the Postgres container
migrate-db:
	docker compose exec -T db psql -U ${DB_USER} -d ${DB_NAME} < migrations/init_schema.up.sql
	docker compose exec -T db psql -U ${DB_USER} -d ${DB_NAME} < migrations/sample_data.up.sql

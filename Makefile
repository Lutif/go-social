# Load environment variables from .env file
ifneq (,$(wildcard .env))
    include .env
    export
endif

MIGRATIONS_PATH = ./cmd/migrate/migrations


.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Using DATABASE_URL: $(DATABASE_URL)"
	@migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATIONS_PATH) -database $(DATABASE_URL) down $(filter-out $@,$(MAKECMDGOALS))

#server
.PHONY: run-server
run-server:
	@go run cmd/api/*.go

.PHONY: watch-server
watch-server:
	@air

.PHONY: db-seed

db-seed:
	@echo "Seeding database with test data using scripts/db-seed.sql"
	@psql "postgres://admin:password@localhost/social?sslmode=disable" -f scripts/db-seed.sql
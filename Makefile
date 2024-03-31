# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	
	@go build -o main cmd/api/api.go

# Run the application
run:
	@go run cmd/api/api.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean

# Migrations
backup-dev:
	@pg_dump $$PG_DEV_URL -W -Ft > ./migrations/backup/backup_`date +%s`.dev.dump

backup-prod:
	@pg_dump $$PG_PROD_URL -W -Ft > ./migrations/backup/backup_`date +%s`.prod.dump

goose-up:
	@echo "Migrating database..."
	@goose -dir ./migrations postgres "host=$$PG_HOST user=$$PG_USER dbname=$$PG_DATABASE password=$$PG_PASSWORD port=$$PG_PORT sslmode=disable" up
	@goose -dir ./migrations postgres $$PG_DEV_URL up
	@goose -dir ./migrations postgres $$PG_PROD_URL up
	@echo "Migration complete"

goose-down:
	@echo "Migrating database..."
	@goose -dir ./migrations postgres "host=$$PG_HOST user=$$PG_USER dbname=$$PG_DATABASE password=$$PG_PASSWORD port=$$PG_PORT sslmode=disable" down
	@goose -dir ./migrations postgres $$PG_DEV_URL down
	@goose -dir ./migrations postgres $$PG_PROD_URL down
	@echo "Migration complete"

goose-status:
	@echo "Local Status of Migrations"
	@goose -dir ./migrations postgres "host=$$PG_HOST user=$$PG_USER dbname=$$PG_DATABASE password=$$PG_PASSWORD port=$$PG_PORT sslmode=disable" status
	@echo "Develop Status of Migrations"
	@goose -dir ./migrations postgres $$PG_DEV_URL status
	@echo "Production Status of Migrations"
	@goose -dir ./migrations postgres $$PG_PROD_URL status


goose-create:
	@read -p "Enter migration name: " NAME; \
  goose -dir ./migrations create $$NAME sql
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

goose-up:
	@goose -dir ./migrations postgres "host=$$DB_HOST user=$$DB_USERNAME dbname=$$DB_DATABASE password=$$DB_PASSWORD port=$$DB_PORT sslmode=disable" up

goose-down:
	@goose -dir ./migrations postgres "host=$$DB_HOST user=$$DB_USERNAME dbname=$$DB_DATABASE password=$$DB_PASSWORD port=$$DB_PORT sslmode=disable" down

goose-status:
	@goose -dir ./migrations postgres "host=$$DB_HOST user=$$DB_USERNAME dbname=$$DB_DATABASE password=$$DB_PASSWORD port=$$DB_PORT sslmode=disable" status

goose-create:
	@read -p "Enter migration name: " NAME; \
  goose -dir ./migrations create $$NAME sql
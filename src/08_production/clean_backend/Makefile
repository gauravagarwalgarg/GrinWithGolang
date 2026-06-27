# Makefile for the clean backend project
#
# Usage:
#   make run        - Run the server locally
#   make build      - Build the binary
#   make test       - Run all tests
#   make docker     - Build and start with Docker Compose
#   make clean      - Remove build artifacts

.PHONY: run build test docker docker-down clean tidy

# Run the server locally (requires MongoDB running)
run:
	go run cmd/main.go

# Build the binary
build:
	go build -o bin/server cmd/main.go

# Run all tests
test:
	go test ./... -v

# Download and tidy dependencies
tidy:
	go mod tidy

# Docker Compose (production-like)
docker:
	docker-compose up -d --build

# Docker Compose (local dev with hot-reload)
docker-local:
	docker-compose -f docker-compose.local.yaml up --build

# Docker Compose down (all)
docker-down:
	docker-compose down
	docker-compose -f docker-compose.local.yaml down

# Clean build artifacts
clean:
	rm -rf bin/

.PHONY: all build run test clean docker-build docker-run help

all: build

build:
	go build -o bin/pack-calculator ./cmd/server

run: build
	./bin/pack-calculator

run-custom: build
	./bin/pack-calculator -pack-sizes="23,31,53"

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

benchmark:
	go test -bench=. -benchmem ./internal/calculator

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

docker-build:
	docker build -t pack-calculator:latest .

docker-run: docker-build
	docker run -p 8080:8080 pack-calculator:latest

docker-compose-up:
	docker-compose up --build

docker-compose-down:
	docker-compose down

fmt:
	go fmt ./...

lint:
	golangci-lint run

help:
	@echo "Available targets:"
	@echo "  make build              - Build the application"
	@echo "  make run                - Run the application locally"
	@echo "  make run-custom         - Run with custom pack sizes (23,31,53)"
	@echo "  make test               - Run unit tests"
	@echo "  make test-coverage      - Run tests with coverage report"
	@echo "  make benchmark          - Run benchmarks"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Build and run Docker container"
	@echo "  make docker-compose-up  - Start with docker-compose"
	@echo "  make docker-compose-down- Stop docker-compose services"
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Lint code"
	@echo "  make help               - Show this help"

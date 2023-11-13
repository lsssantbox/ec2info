.PHONY: help deps lint test build run clean

.DEFAULT_GOAL := help

# Help section
help:
	@echo "Usage:"
	@echo "  make deps       - Download dependencies"
	@echo "  make lint       - Run linters"
	@echo "  make test       - Run unit tests"
	@echo "  make build      - Build the application"
	@echo "  make run        - Build and run the application"
	@echo "  make clean      - Remove built binaries"

# Download dependencies
deps:
	@go get -v -t -d ./...

# Run linters
lint:
	@golangci-lint run

# Run unit tests
test: deps
	@go test -v ./...

# Build the application
build: clean
	@go build -o ec2Info

# Build and run the application
run: deps build
	@./ec2Info

# Remove built binaries
clean:
	@rm -f ec2Info

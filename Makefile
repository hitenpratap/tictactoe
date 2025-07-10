# Makefile for the Go Tic-Tac-Toe Project

# --- Variables ---
# Set the name of the application binary
BINARY_NAME=tictactoe
# Set the name of the Docker Compose service
DOCKER_SERVICE_NAME=tictactoe

# --- Go Commands ---
# Default Go command
GO=go
# Go build command
GO_BUILD=$(GO) build
# Go run command
GO_RUN=$(GO) run
# Go test command
GO_TEST=$(GO) test
# Go clean command
GO_CLEAN=$(GO) clean

# --- Docker Commands ---
# Default Docker Compose command
DOCKER_COMPOSE=docker-compose

.PHONY: all build run test clean docker-build docker-run docker-clean help

# The default target executed when you just run `make`
all: build

# Build the Go application binary
build:
	@echo "Building the application..."
	$(GO_BUILD) -o $(BINARY_NAME) .

# Run the application locally
run:
	@echo "Running the application locally..."
	$(GO_RUN) .

# Run the unit tests
test:
	@echo "Running tests..."
	$(GO_TEST) -v ./...

# Clean up the built binary
clean:
	@echo "Cleaning up..."
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)

# Build the Docker image using docker-compose
docker-build:
	@echo "Building Docker image..."
	$(DOCKER_COMPOSE) build

# Run the application inside a Docker container
docker-run:
	@echo "Running application in Docker..."
	$(DOCKER_COMPOSE) up --build

# Stop and remove all Docker containers, images, and volumes for this service
docker-clean:
	@echo "Tearing down Docker environment (containers, images, volumes)..."
	$(DOCKER_COMPOSE) down --rmi all -v

# Display a help message with all available commands
help:
	@echo "Available commands:"
	@echo "  make build          - Compiles the Go application"
	@echo "  make run            - Runs the application locally"
	@echo "  make test           - Runs the unit tests"
	@echo "  make clean          - Removes the compiled binary"
	@echo "  make docker-build   - Builds the Docker image"
	@echo "  make docker-run     - Builds and runs the application in Docker"
	@echo "  make docker-clean   - Removes all Docker assets for this project"
	@echo "  make help           - Shows this help message"


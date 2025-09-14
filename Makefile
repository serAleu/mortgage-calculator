PROJECT_NAME    := mortgage-calculator
DOCKER_USERNAME := seraleu
DOCKER_IMAGE    := $(DOCKER_USERNAME)/$(PROJECT_NAME)
GO_PACKAGES     := ./...
TEST_FLAGS      := -v -race -cover -coverprofile=coverage.out
BIN_DIR         := bin
BINARY_NAME     := server
PORT            := 8080


.PHONY: all
all: build test lint ## Build, test and lint the project

.PHONY: build
build: ## Build the binary
	@echo "🚀 Building binary..."
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/server

.PHONY: run
run: ## Run the application locally
	@echo "🚀 Starting server on port $(PORT)..."
	go run ./cmd/server

.PHONY: stop
stop: ## Stop the local application (for processes started with 'make run')
	@echo "🛑 Stopping local server..."
	@pkill -f "$(BINARY_NAME)" || true

.PHONY: clean
clean: ## Clean build artifacts and test cache
	@echo "🧹 Cleaning up..."
	go clean
	rm -rf $(BIN_DIR) coverage.out
	go clean -testcache

.PHONY: test
test: ## Run tests with coverage and race detection
	@echo "🧪 Running tests with coverage..."
	go test $(TEST_FLAGS) $(GO_PACKAGES)

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "🧪 Running tests in verbose mode..."
	go test -v $(GO_PACKAGES)

.PHONY: test-coverage
test-coverage: test ## Run tests and show coverage report
	@echo "📊 Generating coverage report..."
	go tool cover -html=coverage.out

.PHONY: test-package
test-package: ## Run tests for specific package (usage: make test-package pkg=./internal/controller)
	@echo "🧪 Running tests for package: $(pkg)"
	go test $(TEST_FLAGS) $(pkg)

.PHONY: lint
lint: ## Run golangci-lint
	@echo "🔍 Running linter..."
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "🔍 Running linter with fixes..."
	golangci-lint run --fix

.PHONY: fmt
fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	go fmt $(GO_PACKAGES)

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE):latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Starting Docker container..."
	docker run -d -p $(PORT):$(PORT) --name $(PROJECT_NAME) $(DOCKER_IMAGE):latest

.PHONY: docker-stop
docker-stop: ## Stop Docker container
	@echo "🛑 Stopping Docker container..."
	docker stop $(PROJECT_NAME) || true

.PHONY: docker-rm
docker-rm: ## Remove Docker container
	@echo "🗑️ Removing Docker container..."
	docker rm $(PROJECT_NAME) || true

.PHONY: docker-clean
docker-clean: docker-stop docker-rm ## Stop and remove Docker container

.PHONY: docker-logs
docker-logs: ## Show Docker container logs
	docker logs -f $(PROJECT_NAME)

.PHONY: docker-push
docker-push: docker-build ## Build and push Docker image to registry
	@echo "📤 Pushing Docker image to registry..."
	docker push $(DOCKER_IMAGE):latest

.PHONY: docker-pull
docker-pull: ## Pull Docker image from registry
	@echo "📥 Pulling Docker image from registry..."
	docker pull $(DOCKER_IMAGE):latest

.PHONY: docker-shell
docker-shell: ## Open shell in Docker container
	docker exec -it $(PROJECT_NAME) /bin/sh

.PHONY: deps
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	go mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "🔄 Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: deps-vendor
deps-vendor: ## Vendor dependencies
	@echo "📦 Vendoring dependencies..."
	go mod vendor

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version: ## Show Go version
	@go version

.DEFAULT_GOAL := help
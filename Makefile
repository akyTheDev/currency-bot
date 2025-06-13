# Name of the binary produced by 'build'
BINARY := currency-bot

# Directory containing goose migrations
MIGRATIONS_DIR := migrations

# Default database URL (can be overridden in environment or by passing on make command line)
DATABASE_URL ?= $DATABASE_URL

# Go commands
GO := go

# Go flags
GOFLAGS :=

# Where to place built binaries
BIN_DIR := bin

# Entrypoint package for 'build' and 'run'
MAIN_PKG := ./cmd/$(BINARY)

# ----------------------------------------
# Phony targets
# ----------------------------------------
.PHONY: all deps fmt vet test build run migrate-up migrate-down clean cover

# Default target: show help
all:
	@echo "Usage:"
	@echo "  make deps          # Download Go module dependencies"
	@echo "  make fmt           # Run go fmt on all packages"
	@echo "  make vet           # Run go vet on all packages"
	@echo "  make test          # Run all tests"
	@echo "  make cover         # Run tests with coverage report"
	@echo "  make build         # Build the binary"
	@echo "  make run           # Run the application (uses DATABASE_URL from env)"
	@echo "  make migrate-up    # Apply goose migrations (requires goose installed)"
	@echo "  make migrate-down  # Roll back goose migrations"
	@echo "  make clean         # Remove built artifacts"

# ----------------------------------------
# Download dependencies
# ----------------------------------------
deps:
	@echo "==> Downloading dependencies..."
	$(GO) mod download

# ----------------------------------------
# Formatting
# ----------------------------------------
fmt:
	@echo "==> Running go fmt..."
	$(GO) fmt ./...

# ----------------------------------------
# Vetting
# ----------------------------------------
vet:
	@echo "==> Running go vet..."
	$(GO) vet ./...

# ----------------------------------------
# Testing
# ----------------------------------------
test:
	@echo "==> Running all tests..."
	$(GO) test $(GOFLAGS) ./...

cover:
	@echo "==> Running tests with coverage..."
	$(GO) test $(GOFLAGS) -coverprofile=coverage.out ./...
	@echo "Coverage report written to coverage.out"

# ----------------------------------------
# Build & Run
# ----------------------------------------
build: fmt vet test
	@echo "==> Building binary: $(BINARY)"
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY) $(MAIN_PKG)
	@echo "Binary created: $(BIN_DIR)/$(BINARY)"

run:
	@echo "==> Running $(BINARY)..."
	# Ensure DATABASE_URL is set in the environment before running
	if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL is not set."; \
		exit 1; \
	fi
	$(GO) run $(MAIN_PKG)

# ----------------------------------------
# Database migrations (requires goose)
# ----------------------------------------
migrate-up:
	@echo "==> Applying migrations (up) in $(MIGRATIONS_DIR)..."
	@echo "Using DATABASE_URL=$(DATABASE_URL)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	@echo "==> Rolling back last migration in $(MIGRATIONS_DIR)..."
	@echo "Using DATABASE_URL=$(DATABASE_URL)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

# ----------------------------------------
# Clean up
# ----------------------------------------
clean:
	@echo "==> Cleaning up binaries and temp files..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out


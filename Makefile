# Use git describe to get a version string.
# Example: v1.0.0-3-g1234567
# Fallback to 'dev' if not in a git repository.
VERSION ?= $(shell git describe --tags --always --dirty --first-parent 2>/dev/null || echo "dev")

# Go parameters
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
GO_CLEAN=$(GO_CMD) clean
GO_INSTALL=$(GO_CMD) install

# Binary name
BINARY_NAME=blogserve
NIX_BUILD=result

# Build flags
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

# Frontend parameters
BUN_CMD=bun
FRONTEND_DIR=frontend

.PHONY: all build run fmt clean install build-frontend

all: build

build-frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && $(BUN_CMD) install && $(BUN_CMD) run build

build: build-frontend
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	CGO_ENABLED=0 $(GO_BUILD) $(LDFLAGS) -o $(BINARY_NAME) .

run:
	$(GO_RUN) . --

fmt:
	@echo "Formatting code..."
	cd $(FRONTEND_DIR) && $(BUN_CMD) run format
	$(GO_CMD) fmt ./...

lint:
	cd $(FRONTEND_DIR) && $(BUN_CMD) run lint
	golangci-lint run ./...

clean:
	@echo "Cleaning..."
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(NIX_BUILD)
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/.svelte-kit

install: build
	@echo "Installing $(BINARY_NAME) to $(shell $(GO_CMD) env GOPATH)/bin..."
	$(GO_INSTALL) .

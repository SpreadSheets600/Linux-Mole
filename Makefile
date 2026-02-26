# Makefile for Mole

.PHONY: all build clean release

# Output directory
BIN_DIR := bin

# Binaries
ANALYZE := analyze
STATUS := status

# Source directories
ANALYZE_SRC := ./cmd/analyze
STATUS_SRC := ./cmd/status

# Build flags
LDFLAGS := -s -w

all: build

# Local build (current architecture)
build:
	@echo "Building for local architecture..."
	@if [ "$$(go env GOOS)" = "darwin" ] || [ "$$(go env GOOS)" = "linux" ]; then \
		go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(ANALYZE)-go $(ANALYZE_SRC); \
	else \
		echo "Skipping analyze-go on $$(go env GOOS) (supported: darwin, linux)"; \
	fi
	go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(STATUS)-go $(STATUS_SRC)

# Release build targets (run on native architectures for CGO support)
release-amd64:
	@echo "Building release binaries (amd64)..."
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(ANALYZE)-darwin-amd64 $(ANALYZE_SRC)
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(STATUS)-darwin-amd64 $(STATUS_SRC)

release-arm64:
	@echo "Building release binaries (arm64)..."
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(ANALYZE)-darwin-arm64 $(ANALYZE_SRC)
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(STATUS)-darwin-arm64 $(STATUS_SRC)

release-linux-amd64:
	@echo "Building Linux release binaries (amd64)..."
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(ANALYZE)-linux-amd64 $(ANALYZE_SRC)
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(STATUS)-linux-amd64 $(STATUS_SRC)

release-linux-arm64:
	@echo "Building Linux release binaries (arm64)..."
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(ANALYZE)-linux-arm64 $(ANALYZE_SRC)
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(STATUS)-linux-arm64 $(STATUS_SRC)

clean:
	@echo "Cleaning binaries..."
	rm -f $(BIN_DIR)/$(ANALYZE)-* $(BIN_DIR)/$(STATUS)-* $(BIN_DIR)/$(ANALYZE)-go $(BIN_DIR)/$(STATUS)-go

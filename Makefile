BINARY  := ent
MODULE  := github.com/AWDDude/ent
CMD_PKG := $(MODULE)/cmd

VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Ldflags are only used by 'make build' (single-target dev build).
# Multi-platform release builds are owned by GoReleaser (.goreleaser.yaml).
LDFLAGS := -s -w \
	-X $(CMD_PKG).Version=$(VERSION) \
	-X $(CMD_PKG).Commit=$(COMMIT) \
	-X $(CMD_PKG).BuildDate=$(BUILD_DATE)

GO      := go
GOFLAGS := -trimpath

.DEFAULT_GOAL := build

# ── Build ──────────────────────────────────────────────────────────────────────

.PHONY: build
build: ## Build the binary for the current platform
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) .

.PHONY: build-all
build-all: ## Build for all supported platforms into dist/ (via goreleaser)
	goreleaser build --snapshot --clean

# ── Test ───────────────────────────────────────────────────────────────────────

.PHONY: test
test: ## Run all tests (verbose, with race detector)
	$(GO) test -v -race ./...

.PHONY: test-cover
test-cover: ## Run tests and open HTML coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# ── Quality ────────────────────────────────────────────────────────────────────

.PHONY: lint
lint: ## Run go vet
	$(GO) vet ./...

.PHONY: tidy
tidy: ## Tidy go.mod and go.sum
	$(GO) mod tidy

.PHONY: check
check: tidy lint test ## Run tidy + lint + tests (CI gate)

# ── Clean ──────────────────────────────────────────────────────────────────────

.PHONY: clean
clean: ## Remove build artefacts
	$(GO) clean
	rm -f $(BINARY) coverage.out
	rm -rf $(OUT_DIR)

# ── Help ───────────────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

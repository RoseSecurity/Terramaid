BINARY_NAME=terramaid
VERSION=v1
GO=go

default: help

help: ## List Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build

fmt: ## Format Go files
	gofumpt -w .

build: ## Build Terramaid
	$(GO) build -ldflags="-s -w" -o build/$(BINARY_NAME) main.go

install: ## Install dependencies
	$(GO) install ./...@latest

clean: ## Clean up build artifacts
	$(GO) clean
	rm ./build/$(BINARY_NAME)

run: build ## Run Terramaid
	./build/$(BINARY_NAME)

docs: build ## Generate documentation
	./build/$(BINARY_NAME) docs
	find ./docs -name '*.md' -print0 | xargs -0 sed -i 's/```terrmaid/```go/g'

mkdocs: ## render mkdcos locally
	mkdocs serve

.PHONY: all build install clean run fmt help mkdocs

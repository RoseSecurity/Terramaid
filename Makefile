BINARY_NAME=terramaid
VERSION=v1
GO=go

default: help

help: ## list makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build

fmt: ## format go files
	gofumpt -w .

build: ## build terramaid
	$(GO) build -ldflags="-s -w" -o build/$(BINARY_NAME) main.go

install: ## install deps
	$(GO) install ./...@latest

clean: ## clean up build artifacts
	$(GO) clean
	rm ./build/$(BINARY_NAME)

run: build ## run 
	./build/$(BINARY_NAME)

docs: build ## docs gen
	./build/$(BINARY_NAME) docs

.PHONY: all build install clean run fmt help
BINARY_NAME=terramaid
VERSION=v1
GO=go

all: build

build:
	$(GO) build -ldflags="-s -w" -o build/$(BINARY_NAME) main.go

install:
	$(GO) install ./...@latest

clean:
	$(GO) clean
	rm ./build/$(BINARY_NAME)

run: build
	./build/$(BINARY_NAME)

docs: build
	./build/$(BINARY_NAME) docs

.PHONY: all build install clean run
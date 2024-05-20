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
	rm $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

.PHONY: all build install clean run
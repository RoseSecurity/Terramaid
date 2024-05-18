get:
	go get

build:
		go build -o build/terramaid terramaid.go

deps:
	go mod download

.PHONY: get build deps

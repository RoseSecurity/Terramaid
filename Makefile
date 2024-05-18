get:
	go get

build: get
	env $(if $(GOOS),GOOS=$(GOOS)) $(if $(GOARCH),GOARCH=$(GOARCH)) go build -o build/terramaid main.go

deps:
	go mod download

.PHONY: get build deps

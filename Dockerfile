FROM golang:1.22

WORKDIR /usr/src/terramaid

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/terramaid main.go

ENTRYPOINT ["/usr/local/bin/terramaid"]
CMD ["default_arg"]
LINTER_VERSION=1.30.0

all: lint test build

lint:
	golangci-lint run

test:
	go test -v ./...

build:
	CGO_ENABLED=0 go build -o bin/shrtnr main.go

LINTER_VERSION=1.33.0

.PHONY:all
all: lint test build

.PHONY:frontend
frontend:
	cd frontend && \
	yarn build --mode development && \
	cd -

.PHONY:lint
lint:
	golangci-lint run

.PHONY:test
test:
	go test -v ./...

.PHONY:build
build:
	rice embed-go
	CGO_ENABLED=0 go build -o bin/shrtnr main.go
	rm rice-box.go

.PHONY:docker
docker:
	docker build . -t indiependente/shrtnr


.PHONY:deps
deps:
	go mod download
	go get github.com/GeertJohan/go.rice
	go get github.com/GeertJohan/go.rice/rice

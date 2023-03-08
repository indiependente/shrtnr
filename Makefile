.PHONY:all
all: ui lint test build

.PHONY:ui
ui:
	cd ui && \
	npm run build && \
	cd -

.PHONY:lint
lint:
	golangci-lint run

.PHONY:test
test:
	go test -v -race ./...

.PHONY:build
build:
	CGO_ENABLED=0 go build -o bin/shrtnr .

.PHONY:docker
docker:
	docker build --no-cache -t indiependente/shrtnr -f Dockerfile .

.PHONY:deps
deps:
	go mod download

## Starts the service locally ( all required components )
.PHONY:start
start: docker
	docker compose -p shrtnr \
	up -d --build --no-recreate --remove-orphans

## Stops the running local service and all its dependencies
.PHONY:stop
stop:
	@ docker compose -p shrtnr down --remove-orphans

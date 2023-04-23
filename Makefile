# Makefile

GO_VERSION=1.20.2
POSTGRES_VERSION := 14
version=$(shell git rev-parse --short HEAD)
image := perebaj/ondehj:$(version)

## Build ondehoje service
.PHONY: ondehoje
ondehoje:
	go build -o ./cmd/ondehoje ./cmd/ondehoje

## Build image service
.PHONY: image
image:
	docker build . \
	--build-arg GO_VERSION=$(GO_VERSION) \
	-t ${image} 

## Run ondehoje service
.PHONY: run
run:
	docker run --rm -p 80:8000 ${image}

## Start containers
.PHONY: dev/start
dev/start:
	POSTGRES_VERSION=${POSTGRES_VERSION} GO_VERSION=${GO_VERSION} docker-compose up -d

## Start containers
.PHONY: dev/restart
dev/restart:
	POSTGRES_VERSION=${POSTGRES_VERSION} GO_VERSION=${GO_VERSION} docker-compose restart

.PHONY: dev/logs
dev/logs:
	docker-compose logs -f app

## Stop and remove containers
.PHONY: dev/stop
dev/stop:
	docker-compose down

## Create tables
.PHONY: dev/migrate
dev/migrate:
	go run cmd/migration/main.go

## Display help for all targets
.PHONY: help
help:
	@awk '/^.PHONY: / { \
		msg = match(lastLine, /^## /); \
			if (msg) { \
				cmd = substr($$0, 9, 100); \
				msg = substr(lastLine, 4, 1000); \
				printf "  ${GREEN}%-30s${RESET} %s\n", cmd, msg; \
			} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
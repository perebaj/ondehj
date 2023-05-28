# Makefile

GO_VERSION=1.20.2
POSTGRES_VERSION := 14
GOLANGCI_LINT_VERSION=v1.51.2

version=$(shell git rev-parse --short HEAD)
image := registry.heroku.com/ondehoje/web:$(version)
devimage :=ondehoje-dev

# To avoid downloading deps everytime it runs on containers
gopkg=$(devimage)-gopkg
gocache=$(devimage)-gocache
devrun=docker run --rm \
	-v `pwd`:/app \
	-v $(gopkg):/go/pkg \
	-v $(gocache):/root/.cache/go-build \
	$(devimage)

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

HEROKU_TOKEN=$(shell heroku auth:token)

## Push image service
.PHONY: publish
publish:
	docker login --username=perebaj@gmail.com --password=${HEROKU_TOKEN} registry.heroku.com
	docker push ${image}


## Run ondehoje service
.PHONY: app
app:
	go run cmd/ondehoje/main.go

## Deploy ondehoje service on heroku
.PHONY: heroku/release
heroku/release:
	heroku container:release web --app ondehoje

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


## Stop and remove containers
.PHONY: dev/stop
dev/stop:
	docker-compose down

## Create tables
.PHONY: dev/migrate
dev/migrate:
	go run cmd/migration/main.go

## Create the dev container image
.PHONY: dev/image
dev/image:
	docker build \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION) \
		-t $(devimage) \
		-f Dockerfile.dev \
		.

## Run tests
.PHONY: test
test: dev/image
	$(devrun) go test ./... -cover

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
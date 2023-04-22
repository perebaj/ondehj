# Makefile


POSTGRES_VERSION := 14

## Start containers
.PHONY: dev/start
dev/start:
	POSTGRES_VERSION=${POSTGRES_VERSION} docker-compose up -d

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
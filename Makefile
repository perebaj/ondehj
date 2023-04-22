# Makefile


POSTGRES_VERSION := 14

.PHONY: dev/start
dev/start:
	POSTGRES_VERSION=${POSTGRES_VERSION} docker-compose up -d

.PHONY: dev/stop
dev/stop:
	docker-compose down

.PHONY: dev/migrate
dev/migrate:
	go run cmd/migration/main.go
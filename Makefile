PHONY:
SILENT:
include .env.prod

MIGRATION_NAME ?= new_migration

DB_CONN_DEV = "host=localhost user=postgres password=postgres port=5441 dbname=order sslmode=disable"

FOLDER_PG= migrations/pg

swag:
	swag init -g cmd/main/main.go

build:
	go build -o ./.bin/main ./cmd/main/main.go
run: build
	./.bin/main

migrations-up:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_DEV)   up

migrations-down:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_DEV)   down


migrations-status:
	goose -dir $(FOLDER_PG) postgres $(DB_CONN_DEV)  status

migrations-new:
	goose -dir $(FOLDER_PG) create $(MIGRATION_NAME) sql

migrations-up-prod:
	goose -dir $(FOLDER_PG) postgres "$(DB_CONN_PROD)" up

migrations-down-prod:
	goose -dir $(FOLDER_PG) postgres "$(DB_CONN_PROD)" down

migrations-status-prod:
	goose -dir $(FOLDER_PG) postgres "$(DB_CONN_PROD)" status

docker-build:
	docker build -t candles .

docker:
	docker-compose  up -d

docker-app: docker-build docker


compose-prod:
	docker compose -f docker-compose.prod.yaml up --build -d
testing:
	go test -v -count=1  ./test/...


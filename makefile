# Simple Makefile
.PHONY: run build test clean

# Load environment variables from .env file
include .env
export

run:
	go run cmd/server/main.go

build:
	go build -o build/server cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

migrate-up:
	goose -dir internal/db/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

migrate-down:
	goose -dir internal/db/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

sqlc:
	sqlc generate

dev: sqlc
	air

setup:
	go mod download
	go install github.com/cosmtrek/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
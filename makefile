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
	goose -dir db/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down:
	goose -dir db/migrations postgres "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

sqlc:
	sqlc generate

dev: sqlc
	air

setup:
	go mod download
	go install github.com/cosmtrek/air@v1.62.0
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
include .env

.PHONY: help build run dev lint migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo " make build          - Build the application"
	@echo " make run            - Run the application"
	@echo " make dev            - Run the application in development mode"
	@echo " make lint           - Run linter on the codebase"
	@echo " make migrate-up     - Apply database migrations"
	@echo " make migrate-up     - Rollback database migrations"

build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5445/ecommerce_api?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5445/ecommerce_api?sslmode=disable" down

migrate-force-zero:
	migrate -path db/migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:5445/ecommerce_api?sslmode=disable" force


docker-up:
	docker compose -f docker/docker-compose.yaml --env-file .env up -d

docker-down:
	docker compose -f docker/docker-compose.yaml --env-file .env down

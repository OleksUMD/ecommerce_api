include .env

.PHONY: help build run dev lint migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo " make build                    - Build the application"
	@echo " make run                      - Run the application"
	@echo " make dev                      - Run the application in development mode"
	@echo " make lint                     - Run linter on the codebase"
	@echo " make migrate-up               - Apply database migrations"
	@echo " make migrate-down             - Rollback database migrations"
	@echo " make docs-generate            - Generate Swagger API documentation"
	@echo " make docker-up                - Start project in docker compose"
	@echo " make docker-down              - Stop project in docker compose"
	@echo " make docker-build-app         - Build app docker image"
	@echo " make docker-build-notifier    - Build notifier docker image"

build:
	@echo "Building all binaries"
	@mkdir -p bin
	@for cmd in cmd/*/; do \
		if [ -d "$$cmd" ]; then \
		  binary=$$(basename $$cmd); \
			echo "Building $$binary..."; \
			go build -o bin/$$binary ./$$cmd; \
		fi \
	done

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
	docker compose -f docker/docker-compose.yaml --env-file .docker.env up -d

docker-down:
	docker compose -f docker/docker-compose.yaml --env-file .docker.env down

docker-build-app:
	docker compose -f docker/docker-compose.yaml --env-file .docker.env build app --no-cache

docker-build-notifier:
	docker compose -f docker/docker-compose.yaml --env-file .docker.env build notifier --no-cache

docs-generate:
	mkdir -p docs
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --exclude .git,docs,docker,db

DEFAULT_GOAL := up-watch
ENV_FILE := .env
include $(ENV_FILE)
export

build:
	@echo "Building..."
	@go build -o ./bin/main cmd/api/main.go

swag:
	swag init -g ./cmd/api/main.go --parseDependency -o ./docs/swagger/

run: 
	@go run cmd/api/main.go
	 
run-worker:
	@go run cmd/worker/main.go

sqlc:
	sqlc generate

reset-db: goose-down goose-up

goose-create:
	goose -s create $(name) -dir ./internal/database/migrations sql

goose-up:
	goose up

goose-status:
	goose status


goose-down-to:
	goose down-to $(version)

goose-down:
	goose down-to 0

up:
	docker compose up -d --build

up-watch:
	docker compose up --build --watch 

down:
	docker compose down -v

test:
	go test -v ./...

watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi


.PHONY: all build test clean watch docker-run docker-down itest run run-prod watch docker-up docker-down up down sqlc up-watch

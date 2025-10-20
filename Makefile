# Database configuration
DB_USER=root
DB_PASSWORD=1111
DB_HOST=localhost
DB_PORT=3306
DB_NAME=stk_menu_system
DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# Application commands
.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: build
build:
	go build -o bin/api cmd/api/main.go

.PHONY: dev
dev:
	air

# Database commands
.PHONY: db-create
db-create:
	mysql -u $(DB_USER) -p$(DB_PASSWORD) -e "CREATE DATABASE IF NOT EXISTS $(DB_NAME);"

.PHONY: db-drop
db-drop:
	mysql -u $(DB_USER) -p$(DB_PASSWORD) -e "DROP DATABASE IF EXISTS $(DB_NAME);"

# Migration commands
.PHONY: migrate-up
migrate-up:
	migrate -path database/migrations -database "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	migrate -path database/migrations -database "$(DB_URL)" down

.PHONY: migrate-down-all
migrate-down-all:
	migrate -path database/migrations -database "$(DB_URL)" down -all

.PHONY: migrate-force
migrate-force:
	migrate -path database/migrations -database "$(DB_URL)" force $(VERSION)

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir database/migrations -seq $(name)

# Dependency management
.PHONY: deps
deps:
	go mod download
	go mod tidy

.PHONY: deps-update
deps-update:
	go get -u ./...
	go mod tidy

# Testing
.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Linting
.PHONY: lint
lint:
	golangci-lint run

# Clean
.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out

# Setup (first time)
.PHONY: setup
setup:
	cp .env.example .env
	go mod download
	make db-create
	make migrate-up

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build the application"
	@echo "  make db-create        - Create database"
	@echo "  make db-drop          - Drop database"
	@echo "  make migrate-up       - Run migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-create   - Create new migration (use: make migrate-create name=migration_name)"
	@echo "  make deps             - Download dependencies"
	@echo "  make test             - Run tests"
	@echo "  make setup            - First time setup (create .env, db, and run migrations)"


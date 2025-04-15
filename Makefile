.PHONY: run build migrate-up migrate-down migrate-create compose-up compose-down test

MIGRATIONS_PATH=./migrations
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5436
DB_NAME=auth_db
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

run:
	go run cmd/main.go

build:
	go build -o bin/app cmd/main.go

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down


compose-up:
	docker-compose up --build

compose-down:
	docker-compose down -v

test:
	go test -v -coverprofile=coverage.out ./internal/usecase/...
	go tool cover -func=coverage.out

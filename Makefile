.PHONY: run build migrate-up migrate-down migrate-create compose-up compose-down test

# Переменные для миграций и БД
MIGRATIONS_PATH=./migrations
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5436
DB_NAME=auth_db
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Запуск приложения локально
run:
	go run cmd/main.go

# Сборка бинарника
build:
	go build -o bin/app cmd/main.go

# Применить миграции вверх
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# Откатить миграции вниз
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down


# Запуск docker-compose
compose-up:
	docker-compose up --build

# Остановка docker-compose и удаление volume
compose-down:
	docker-compose down -v

# Тесты с покрытием
test:
	go test -v -coverprofile=coverage.out ./internal/usecase/...
	go tool cover -func=coverage.out

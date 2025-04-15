# Стадия сборки
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из стадии сборки
COPY --from=builder /app/auth-service .

# Копируем .env, если нужно (или монтируем через docker-compose)
COPY .env .env

# Открываем порт (если нужно)
EXPOSE 8080

# Запуск приложения
CMD ["./auth-service"]

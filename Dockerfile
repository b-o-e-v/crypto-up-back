# Stage 1: Сборка приложения
FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app

# Копируем и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем приложение
COPY . .
RUN go version
RUN go build -tags=jsoniter -o crypto-up .

# Stage 2: Финальный образ
FROM alpine:latest

# Устанавливаем необходимые сертификаты
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем исполняемый файл и файл конфигурации из предыдущего этапа
COPY --from=builder /app/crypto-up .
COPY --from=builder /app/.env .

# Устанавливаем команду по умолчанию
CMD ["./crypto-up"]

# docker-compose build --no-cache

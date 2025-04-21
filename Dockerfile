# Строим на основе golang:1.24-alpine
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем все исходные файлы Go в рабочую директорию
COPY . .

# Строим приложение
RUN go build -o app ./cmd

# Минимальный образ для запуска
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный файл из builder
COPY --from=builder /app/app .

# Копируем .env файл
COPY .env .env

# Запускаем приложение
CMD ["./app"]

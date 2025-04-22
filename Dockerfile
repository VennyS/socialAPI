# ─────────────────────────────────────
# 🔨 Сборочный этап
# ─────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

# ─────────────────────────────────────
# 🚀 Минимальный запускной образ с dockerize
# ─────────────────────────────────────
FROM alpine:latest

WORKDIR /root/

# Установка dockerize
ADD https://github.com/jwilder/dockerize/releases/download/v0.6.1/dockerize-linux-amd64-v0.6.1.tar.gz /tmp/
RUN tar -C /usr/local/bin -xzvf /tmp/dockerize-linux-amd64-v0.6.1.tar.gz && \
    rm /tmp/dockerize-linux-amd64-v0.6.1.tar.gz

# Копируем скомпилированное приложение
COPY --from=builder /app/app .

# Копируем переменные окружения
COPY .env .env

# Запуск с ожиданием db и redis
CMD ["dockerize", "-wait", "tcp://db:5432", "-wait", "tcp://redis:6379", "-timeout", "30s", "./app", "-migrate"]
FROM golang:1.21.1-alpine AS builder

WORKDIR /usr/local/src

# Установка необходимых пакетов
RUN apk --no-cache add bash git make gcc gettext musl-dev

# Копирование go.mod и go.sum
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Копирование исходного кода
COPY ./ ./

# Сборка приложения
RUN go build -o ./bin/app cmd/vk-test-task/main.go

# Создание нового образа для запуска приложения
FROM alpine AS runner

# Копирование исполняемого файла из стадии сборки
COPY --from=builder /usr/local/src/bin/app /app
RUN chmod +x /app

# Копирование файла .env в образ
COPY .env .env
COPY config/config_http.yml config/config_http.yml
COPY config/config_db.yml config/config_db.yml
RUN ls -la /app
# Установка рабочей директории
WORKDIR /

# Запуск приложения
CMD ["/app"]
CMD ["/bin/sh", "-c", "/app"]
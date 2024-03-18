FROM golang:1.21.1-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

# Замените apt-get на apk
RUN apk update && apk add postgresql-client

# Сделать wait-for-postgres.sh исполняемым
COPY wait-for-postgres.sh /usr/local/src/wait-for-postgres.sh
RUN chmod +x wait-for-postgres.sh

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app cmd/vk-test-task/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app

COPY .env .env
COPY config/config_http.yml config/config_http.yml
COPY config/config_db.yml config/config_db.yml
RUN ls -la /app
WORKDIR /

CMD ["/app"]

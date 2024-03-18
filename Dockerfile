FROM golang:1.21.1-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app cmd/vk-test-task/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app
RUN chmod +x /app

COPY .env .env
COPY config/config_http.yml config/config_http.yml
COPY config/config_db.yml config/config_db.yml
RUN ls -la /app
WORKDIR /

CMD ["/app"]
CMD ["/bin/sh", "-c", "/app"]
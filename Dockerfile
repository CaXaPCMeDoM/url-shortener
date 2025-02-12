FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrator/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/url-shortener ./cmd/url-shortener/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/url-shortener /app/url-shortener

COPY ./migrations /app/migrations

COPY .env /app/.env
COPY ./config/local.yaml /app/config/local.yaml

RUN echo -e '#!/bin/sh\n\
/app/migrate -migrations-path /app/migrations up\n\
exec /app/url-shortener' > /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

EXPOSE 8090 44044

CMD ["/app/entrypoint.sh"]
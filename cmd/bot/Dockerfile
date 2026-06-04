# Builder
FROM golang:alpine AS builder
WORKDIR /app

COPY . .

RUN go build -o bot ./cmd/bot/main.go
RUN go build -o migrate ./cmd/migrate/main.go

# Result image
FROM alpine
WORKDIR /app

COPY .env .env
COPY sql sql
COPY "cmd/bot/entrypoint.sh" entrypoint.sh
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/bot /app/bot
RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]
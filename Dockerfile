FROM golang:alpine AS builder

WORKDIR /app

RUN go version
ENV GOPATH=/

COPY . .

RUN go build  -o bot ./cmd/bot/main.go

FROM alpine
COPY --from=builder /app/bot /app/bot

CMD ["./bot"]
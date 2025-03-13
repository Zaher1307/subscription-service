FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /subscription-service ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /subscription-service .

EXPOSE 8080

CMD ["./subscription-service"]

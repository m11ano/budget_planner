FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY gateway /app/gateway
COPY auth /app/auth

WORKDIR /app/gateway

RUN go mod download

WORKDIR /app/gateway
RUN go build -o app ./cmd/main.go


FROM golang:1.24-alpine AS release

WORKDIR /app

COPY --from=builder /app/gateway/app .
COPY --from=builder /app/gateway/configs/base.yml /app/configs/base.yml

CMD ["./app"]
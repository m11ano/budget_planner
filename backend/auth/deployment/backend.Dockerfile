FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY auth /app/auth

WORKDIR /app/auth

RUN go mod download

WORKDIR /app/auth
RUN go build -o app ./cmd/main.go


FROM golang:1.24-alpine AS release

WORKDIR /app

COPY --from=builder /app/auth/app .
COPY --from=builder /app/auth/configs/base.yml /app/configs/base.yml
COPY --from=builder /app/auth/migrations /app/migrations 

CMD ["./app"]
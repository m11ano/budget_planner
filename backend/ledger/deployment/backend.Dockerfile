FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ledger /app/ledger
COPY auth /app/auth

WORKDIR /app/ledger

RUN go mod download

WORKDIR /app/ledger
RUN go build -o app ./cmd/main.go


FROM golang:1.24-alpine AS release

WORKDIR /app

COPY --from=builder /app/ledger/app .
COPY --from=builder /app/ledger/configs/base.yml /app/configs/base.yml
COPY --from=builder /app/ledger/migrations /app/migrations

CMD ["./app"]
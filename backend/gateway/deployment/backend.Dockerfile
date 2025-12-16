# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download


COPY . .
ARG TARGETOS TARGETARCH
ARG VERSION="dev"
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -tags timetzdata -buildvcs=false \
    -ldflags="-s -w -X 'main.version=${VERSION}'" \
    -o /out/app ./cmd/backend

FROM gcr.io/distroless/static-debian12:nonroot AS release
WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /out/app ./app
COPY --from=build --chown=nonroot:nonroot --chmod=555 /src/configs ./configs
COPY --from=build --chown=nonroot:nonroot --chmod=555 /src/migrations ./migrations

USER nonroot:nonroot
ENTRYPOINT ["./app"]

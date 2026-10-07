FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go install github.com/Weit145/simple-log/cmd/simple-log@3cdeda5df8ce3fe7e37433d3a98e693573b1e5c1

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app/main.go

FROM golang:1.26-alpine AS goose-builder

RUN CGO_ENABLED=0 GOBIN=/out go install -tags='no_azuresql no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' github.com/pressly/goose/v3/cmd/goose@v3.28.0

FROM alpine:3.22 AS migrator

WORKDIR /app
COPY --from=goose-builder /out/goose /usr/local/bin/goose
COPY migrations ./migrations
ENV GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=/app/migrations
ENTRYPOINT ["goose", "up"]

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/app ./app
COPY --from=builder /go/bin/simple-log /usr/local/bin/simple-log

USER app

EXPOSE 8080

ENTRYPOINT ["simple-log", "--", "./app"]

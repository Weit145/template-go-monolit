FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go install github.com/Weit145/simple-log/cmd/simple-log@3cdeda5df8ce3fe7e37433d3a98e693573b1e5c1

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app/main.go

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

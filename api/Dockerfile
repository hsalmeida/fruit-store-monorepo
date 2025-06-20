# ── builder ───────────────────────────────────────────────────────────────────
FROM golang:1.24.4-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/api cmd/main.go

# ── final image ───────────────────────────────────────────────────────────────
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/api .

EXPOSE 8080

ENTRYPOINT ["./api"]

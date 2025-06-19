# ── builder ───────────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder
WORKDIR /src

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Adjust path to your consumer’s main.go if needed (e.g. cmd/consumer/main.go)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/user-service cmd/consumer/main.go

# ── final image ───────────────────────────────────────────────────────────────
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/user-service .

ENTRYPOINT ["./user-service"]

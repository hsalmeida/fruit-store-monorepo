# ── builder ───────────────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder
WORKDIR /src

# Install git for possible go module fetches
RUN apk add --no-cache git

# Copy Go modules manifests and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy all source and build the API binary
COPY . .
# Adjust the path to your main.go if needed (e.g. cmd/main.go or cmd/server/main.go)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/api cmd/main.go

# ── final image ───────────────────────────────────────────────────────────────
FROM alpine:3.18
# Include CA certificates for TLS (e.g. HTTPS clients)
RUN apk add --no-cache ca-certificates

WORKDIR /app
# Copy the statically linked binary
COPY --from=builder /app/api .

# Expose the port the API listens on
EXPOSE 8080

# Start the API
ENTRYPOINT ["./api"]

# ============================
# BUILD STAGE
# ============================
FROM golang:1.25-alpine AS builder


WORKDIR /app

# Required for go mod download
RUN apk add --no-cache git

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o mma ./cmd/mma


# ============================
# RUNTIME STAGE
# ============================
FROM alpine:3.20

WORKDIR /app

# Certificates (MQTT TLS / HTTPS later)
RUN apk add --no-cache ca-certificates

# Copy binary only
COPY --from=builder /app/mma /app/mma

# --------------------------------
# EXPOSE ONLY WHAT WE SERVE
# --------------------------------
# Modbus TCP (standard)
EXPOSE 502

# REST + Web UI (shared HTTP)
EXPOSE 8080

# Run
ENTRYPOINT ["/app/mma"]

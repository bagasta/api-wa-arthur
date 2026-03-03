# ─── Stage 1: Builder ────────────────────────────────────────────────────────
FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Install gcc for CGO (required by go-sqlite3)
RUN apt-get update && apt-get install -y gcc libc6-dev && rm -rf /var/lib/apt/lists/*

# Download dependencies first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY *.go ./

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o whatsapp-endpoint .

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies (sqlite3 libs + ca-certs for HTTPS)
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

# Copy compiled binary from builder
COPY --from=builder /app/whatsapp-endpoint .

# Volume for WhatsApp session persistence
VOLUME ["/app/data"]

EXPOSE 8200

CMD ["./whatsapp-endpoint"]

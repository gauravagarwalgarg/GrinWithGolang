# Multi-stage build for minimal production image
#
# LEARNING NOTES:
# - Stage 1 (builder): compiles Go binary with all dependencies
# - Stage 2 (runtime): copies only the binary into a minimal Alpine image
# - Result: ~15MB image vs ~800MB if you ship the entire Go toolchain

# --- Build stage ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency files first (Docker layer caching: deps change less often)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# --- Runtime stage ---
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]

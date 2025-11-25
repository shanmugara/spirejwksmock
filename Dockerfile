# Multi-stage Dockerfile for spirejwksmock
# Builder stage
FROM golang:1.25-alpine AS builder

# Install git for `go mod` if needed
RUN apk add --no-cache git ca-certificates

WORKDIR /src

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static, stripped binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags "-s -w" -o /spirejwksmock

# Final stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /spirejwksmock /usr/local/bin/spirejwksmock
RUN chown appuser:appgroup /usr/local/bin/spirejwksmock

USER appuser
ENTRYPOINT ["/usr/local/bin/spirejwksmock"]


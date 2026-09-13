# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and SSL certificates for downloading dependencies
RUN apk add --no-cache git ca-certificates

# Copy source code
COPY . .

# Ensure go.mod uses a standard supported Go version inside Docker
RUN sed -i 's/^go 1\..*/go 1.22/' go.mod

# Download dependencies & compile static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install runtime SSL certificates
RUN apk --no-cache add ca-certificates

# Copy compiled binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

# Default Environment Variables
ENV PORT=8080
ENV JWT_SECRET="super-secret-key-change-in-production"
ENV DB_PATH="tickets.db"

CMD ["./main"]

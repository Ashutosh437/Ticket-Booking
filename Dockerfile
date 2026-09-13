# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy source code
COPY . .

# Align go.mod version and synchronize checksums inside container
RUN sed -i 's/^go 1\..*/go 1.22/' go.mod && go mod tidy

# Build static Linux binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install runtime SSL certificates
RUN apk --no-cache add ca-certificates

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

EXPOSE 8080

# Environment Defaults
ENV PORT=8080
ENV JWT_SECRET="super-secret-key-change-in-production"
ENV DB_PATH="tickets.db"

CMD ["./main"]

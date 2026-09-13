# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency definitions
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build lightweight binary (CGO_ENABLED=0 for pure Go compilation)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for secure HTTP requests
RUN apk --no-cache add ca-certificates

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

EXPOSE 8080

# Environment Defaults
ENV PORT=8080
ENV JWT_SECRET="super-secret-key-change-in-production"
ENV DB_PATH="tickets.db"

CMD ["./main"]

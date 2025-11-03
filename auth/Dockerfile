FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations (optional, if you want to run them from inside container)
COPY --from=builder /app/migrations ./migrations

# Copy .env files (in production, use env vars directly)
COPY --from=builder /app/.env* ./

EXPOSE 8082

CMD ["./main"]

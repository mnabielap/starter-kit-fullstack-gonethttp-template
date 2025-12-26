# Stage 1: Build the Go Application
FROM golang:1.25.4-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (CGO disabled for static linking)
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Stage 2: Create the Production Image
FROM alpine:latest

# Install certificates for HTTPS and bash for entrypoint
RUN apk --no-cache add ca-certificates bash

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy necessary assets (Templates, Static files, Migrations)
# Go templates are parsed at runtime, so we need the physical files
COPY --from=builder /app/web ./web
COPY --from=builder /app/migrations ./migrations

# Copy entrypoint script
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Create a directory for SQLite persistence (if needed)
RUN mkdir -p /app/db

# Expose the application port
EXPOSE 5005

# Set the entrypoint
ENTRYPOINT ["./entrypoint.sh"]
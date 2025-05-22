FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# Use a smaller image for the final application
FROM alpine:latest

WORKDIR /app

# Install PostgreSQL client
RUN apk --no-cache add postgresql-client

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the database schema
COPY --from=builder /app/database/sql_schema.sql ./database/

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"] 
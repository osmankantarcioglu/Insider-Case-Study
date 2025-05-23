FROM golang:1.20-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"] 
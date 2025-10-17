# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Download dependencies and create go.sum
RUN go mod download && go mod tidy

# Copy source code
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Install ca-certificates untuk HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary dari builder
COPY --from=builder /app/main .

# Copy .env file jika ada (optional)
COPY .env* ./

# Expose port
EXPOSE 8080

# Run aplikasi
CMD ["./main"]
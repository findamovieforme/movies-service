# Stage 1: Build the Go app
FROM golang:1.23-alpine AS builder

# Install any build dependencies (optional)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o movies-service main.go

# Stage 2: Create a minimal image with just the binary
FROM alpine:latest

# Install certificates to allow HTTPS connections (important for many Go apps)
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/movies-service /app/movies-service

# Expose port 8080
EXPOSE 8081

# Run the app
CMD ["/app/movies-service"]

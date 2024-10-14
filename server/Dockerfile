# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first (for better layer caching)
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o go_server .

# Command to run your Go app
CMD ["./go_server"]

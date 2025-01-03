# Stage 1: Build the binary
FROM golang:1.23.4 AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files for dependency resolution
COPY go.mod go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the application code
COPY . .

# Build the application
RUN go build -o udm-client ./main.go

# Stage 2: Create a minimal image to run the binary
FROM alpine:latest

# Install required certificates for HTTPS communication
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/udm-client .

# Expose the port for the microservice
EXPOSE 8081

# Command to run the binary
CMD ["./udm-client"]


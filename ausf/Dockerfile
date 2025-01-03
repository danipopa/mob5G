FROM golang:1.23.4 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ausf-service ./main.go

# Use a minimal image for the runtime
FROM alpine:latest

# Install bash
RUN apk apk update && apk upgrade && apk add bash

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory for the runtime stage
WORKDIR /app

# Copy the compiled binary
COPY --from=builder /app/ausf-service .

# Expose the port
EXPOSE 8080

# Start the service
CMD ["./ausf-service"]


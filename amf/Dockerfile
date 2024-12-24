# Base image for Go
FROM golang:1.23.4 as builder

# Set working directory in the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the project
RUN go build -o amf ./cmd/amf

# Use a minimal image for production
FROM alpine:latest

# Set working directory in the production container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/amf .

# Expose SCTP and other required ports
EXPOSE 38412 8080

# Command to run the executable
CMD ["./amf"]


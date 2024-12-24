# Use official Golang image
FROM golang:1.23.4

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project directory into the container
COPY . .

# Build the Go application
RUN go build -o nrf ./cmd/main.go

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./nrf"]


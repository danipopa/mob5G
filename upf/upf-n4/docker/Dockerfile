# Dockerfile for UPF-N4
FROM golang:1.23.4 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o upf-n4 ./cmd/server

FROM debian:bullseye
WORKDIR /app
COPY --from=builder /app/upf-n4 /app/
COPY configs/config.yaml /app/config.yaml
ENTRYPOINT ["/app/upf-n4"]


FROM golang:1.23.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o smf-session-controller ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/smf-session-controller .
EXPOSE 8082
CMD ["./smf-session-controller"]


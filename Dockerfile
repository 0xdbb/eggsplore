# ---------- Builder Stage ----------
FROM golang:1.23.8-alpine3.20 AS builder

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go


# ---------- Runtime Stage ----------
FROM alpine:latest

WORKDIR /app

# Install necessary system dependencies
RUN apk --no-cache add ca-certificates

# Copy only the built binary and .env file
COPY --from=builder /app/main .
# COPY .env .env

# Expose the API port
EXPOSE 8080

# Run the binary
CMD ["./main"]

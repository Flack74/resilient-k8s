FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/chaos-operator ./cmd/chaos-operator

# Create a minimal image
FROM alpine:3.16

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache ca-certificates bash curl

# Copy the binary from the builder stage
COPY --from=builder /app/chaos-operator /app/chaos-operator

# Create entrypoint script
COPY deployments/docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Run the application
ENTRYPOINT ["/app/entrypoint.sh"]
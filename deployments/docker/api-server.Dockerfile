FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api-server ./cmd/api-server

# Create a minimal image
FROM alpine:3.16

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api-server /app/api-server

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["/app/api-server"]
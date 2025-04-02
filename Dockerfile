FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o go_short .

# Use a smaller image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/go_short .

# Copy the .env file
COPY .env .

# Create a non-root user and switch to it
RUN adduser -D -g '' appuser
USER appuser

# Command to run the executable
CMD ["./go_short"]

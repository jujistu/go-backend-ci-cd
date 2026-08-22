####################################################################
# Builder Stage
####################################################################
FROM golang:alpine AS builder

LABEL MAINTAINER="georgetheprogrammer@gmail.com"

WORKDIR /go/src/github.com/ong-gtp/go-chat

# Install git for Go dependencies
RUN apk add --no-cache git

# Copy dependency files first for better Docker layer caching
COPY go.mod go.sum ./

RUN go mod download

# Copy application source
COPY . .

# Build the API service
RUN go build -o gochatapp


####################################################################
# Final Stage
####################################################################
FROM alpine:3.22

# Update Alpine packages
RUN apk upgrade --no-cache

# Create a dedicated non-root user and group
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy application binary
COPY --from=builder /go/src/github.com/ong-gtp/go-chat/gochatapp ./gochatapp

# Ensure the application belongs to the non-root user
RUN chown appuser:appgroup /app/gochatapp

# Run container as non-root user
USER appuser

# Application port
EXPOSE 9010

# Start application
CMD ["./gochatapp"]
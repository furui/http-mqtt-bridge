# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Target architecture arguments supplied automatically by Buildx
ARG TARGETOS
ARG TARGETARCH

# Build the application
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o http-mqtt-bridge main.go

# Stage 2: Create the minimal runtime image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/http-mqtt-bridge .

# Expose the default HTTP port
EXPOSE 8080

# Run the application
CMD ["./http-mqtt-bridge"]

# Stage 1: Build the Go binary
FROM golang:1.23.2 AS builder
ARG CGO_ENABLED=0
ARG VERSION

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o zkill ./main.go

# Stage 2: Create the final image using distroless
FROM gcr.io/distroless/static:latest

# Copy static assets if any (ensure /app/static exists in builder stage)
COPY --from=builder /app/static /static

# Copy the Go binary
COPY --from=builder /app/zkill /zkill

COPY --from=builder /app/.env /.env

# Define a volume for data persistence
VOLUME ["/data"]

# Set the entrypoint to the Go binary
ENTRYPOINT ["/zkill"]

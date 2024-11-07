# Stage 1: Build the Go application
FROM golang:1.22.6 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o download-exporter .

# Stage 2: Create a minimal image
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/download-exporter /download-exporter
COPY --from=builder /app/configs/.env /configs/.env
COPY --from=builder /files.yaml /files.yaml

USER 65532:65532
EXPOSE 8181

# Set the entrypoint to the Go binary
CMD ["/download-exporter"]

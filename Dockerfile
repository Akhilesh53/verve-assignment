# Use the official Golang image for version 1.22.4 as the base image
FROM golang:1.22.4-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o myapp

# Use a minimal base image for the final stage
FROM alpine:latest

# Install necessary packages (e.g., CA certificates for HTTPS)
RUN apk --no-cache add ca-certificates

# Set the working directory for the final image
WORKDIR /root/

# Copy the binary from the builder stage to the final image
COPY --from=builder /app/myapp .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]

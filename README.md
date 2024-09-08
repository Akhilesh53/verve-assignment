# verve-assignment: Unique Request Counter API

## Overview

This Go application is a REST service designed to handle and track unique request IDs. It features:
- **GET Endpoint**: `/api/verve/accept` to accept integer `id` and an optional `endpoint` query parameter.
- **Logging**: Logs the count of unique requests every minute.
- **HTTP Requests**: Can send HTTP GET or POST requests to a specified endpoint with the count of unique requests.
- **Extensions**:
  - **POST Request**: Sends a POST request with a JSON payload.
  - **Load Balancer Compatibility**: Handles ID deduplication across multiple instances.
  - **Distributed Streaming**: Sends the count of unique requests to a Redis instance.

## Requirements

- Go 1.22.4
- Docker (for containerization)
- Redis (for distributed deduplication and streaming)

## Build and Run (On Local)

Buid the application:

```go build -o myapp main.go```

Run the application:

```./myapp```

## Build and Run (Docker)

Build the Docker image:

```docker build -t myapp .```

Run the Docker container:

```docker run -p 8080:8080 myapp```


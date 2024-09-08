# Thought Process

## High-Level Overview

### Objective
The goal of this project is to develop a high-performance REST service using Go that can:
1. Handle at least 10,000 requests per second.
2. Deduplicate requests based on an integer `id`.
3. Send periodic reports of unique request counts to an external endpoint.
4. Ensure deduplication works correctly across multiple instances, even when behind a load balancer.

### Implementation Approach

1. **Core Service**:
   - **REST Endpoint**: Implemented a `GET /api/verve/accept` endpoint using the Gin framework.
   - **Request Handling**: Processes incoming requests by extracting the `id` and optional `endpoint` query parameters. Checks if the ID is unique and performs actions based on that.

2. **Request Deduplication**:
   - **In-Memory Deduplication**: Initially used an in-memory `sync.Map` to track unique IDs within a given minute.
   - **Redis Integration**: Replaced in-memory storage with Redis to handle unique ID tracking across multiple instances. This ensures consistency in environments with load balancers and multiple service instances.

3. **Logging and Reporting**:
   - **Logging**: Uses a standard logger to log the count of unique requests every minute.
   - **External Reporting**: Configured to send unique request counts to a Redis instance for distributed streaming.

4. **HTTP Requests**:
   - **GET Requests**: Sends a GET request with the count of unique requests to an external endpoint.
   - **POST Requests**: Added functionality to send a POST request with a JSON payload containing the count.

5. **Concurrency and Synchronization**:
   - **Locks**: Used to manage concurrent access to Redis and ensure that request counting and logging operations are thread-safe.

6. **Dockerization**:
   - **Dockerfile**: Created a Dockerfile to containerize the Go application, ensuring consistent deployment and easy setup.
   - **Redis Container**: Provides instructions to run a Redis container for development and testing.

### Design Considerations

1. **Scalability**:
   - **Redis**: Leveraged Redis for distributed ID deduplication and count tracking, allowing the service to scale across multiple instances.
   - **Concurrency**: Managed using Go’s synchronization primitives and Redis atomic operations to ensure accuracy and performance.

2. **Fault Tolerance**:
   - **Error Handling**: Implemented robust error handling and logging to capture and address issues during request processing and external communication.

3. **Performance**:
   - **Request Handling**: Optimized the request handling process to achieve high throughput.
   - **Logging**: Efficiently logs unique request counts to avoid performance bottlenecks.

4. **Deployment**:
   - **Docker**: Used Docker for containerization to simplify deployment and ensure the application runs consistently across different environments.
   - **Redis Integration**: Configured Redis for both deduplication and distributed streaming, ensuring seamless integration with the application.

## Repository and Deployment

- **Source Code**: Available at [GitHub Repository](https://github.com/Akhilesh53/verve-assignment.git)
- **Docker Container**: 
  - **Build Command**: `docker build -t myapp .`
  - **Run Command**: `docker run -p 8080:8080 --link redis:redis myapp`
  - **Redis Container**: `docker run -d --name redis -p 6379:6379 redis`

## Conclusion

The implementation leverages modern technologies and best practices to ensure high performance, scalability, and reliability of the REST service. By utilizing Redis for distributed deduplication and Docker for containerization, the solution is designed to handle high request volumes and work seamlessly in various deployment scenarios.


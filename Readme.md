# Reverse Proxy Server in Go

This is a powerful and flexible reverse proxy server written in Go. It provides advanced features such as path-based routing, rate limiting, health checks, and easy configurability. The server allows you to efficiently manage incoming requests, directing them to different backend servers based on the requested path, all while ensuring optimal performance and reliability.

## Example Usage:

1. Install required dependencies using `go mod tidy`.
2. Configure the reverse proxy server using the provided JSON configuration file.

   ```json
   {
     "rateLimiting": {
       "maxRequestsPerMinute": 10,
       "tokenRefillInterval": 600
     },
     "targetServers": {
       "http://backend-server-1": {
         "routePath": "/api/resource1/*",
         "healthPath": "/health"
       },
       "http://backend-server-2": {
         "routePath": "/api/resource2/*",
         "healthPath": "/health"
       }
     },
     "healthCheckFrequency": 5000
   }
3. Run the server using `go run main.go`


## Key Features
- Path-Based Routing:
    - Efficiently forward requests to different backend servers based on the requested path.
- Rate Limiting:
    - Control the number of requests to each backend server using a token bucket rate limiter.
- Health Checks:
    - Periodically check the health of backend servers to ensure reliable routing.
- Configurability:
    - Easily configure the server behavior using a JSON configuration file.
- Easy Integration:
    - Integrate the reverse proxy server into your application effortlessly by following the provided usage instructions.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
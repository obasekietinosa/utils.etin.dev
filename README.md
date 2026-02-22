# utils.etin.dev

Bespoke utilities exposed as endpoints for my use. This project provides a simple, production-quality Go API for various one-off tools.

## Prerequisites

*   [Go](https://go.dev/dl/) 1.24+
*   Make (optional)

## Project Structure

*   `cmd/server`: Application entry point.
*   `internal/handlers`: HTTP handlers for the API endpoints.
*   `internal/middleware`: HTTP middleware (e.g., logging).

## Getting Started

### Running Locally

1.  Clone the repository.
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Run the server:
    ```bash
    make run
    # or
    go run ./cmd/server
    ```

The server will start on port 8080 by default. You can override this by setting the `PORT` environment variable.

### Building

To build the binary:

```bash
make build
# or
go build -o server ./cmd/server
```

## Development

### Running Tests

```bash
make test
# or
go test -v ./...
```

### Linting

```bash
make lint
# or
go vet ./...
```

## API Endpoints

### GET /health

Returns the health status of the service.

**Response:**

```json
{
  "status": "ok"
}
```

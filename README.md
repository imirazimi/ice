# Todo Service - ICE Global Assignment

This project is designed using Clean Architecture to manage Todo items and send each new item to a Redis Stream.

## Getting Started

1. Copy environment variables:

```sh
cp .env.example .env
```

2. Start the backing services:

```sh
make up
```

3. Run database migrations:

```sh
make migrate
```

4. Run the application:

```sh
make run
```

5. Create a new Todo:

Send a POST request to:

```
POST http://localhost:8080/todo
Content-Type: application/json
{
  "description": "test task",
  "dueDate": "2025-01-01T06:00:00Z"
}
```

6. Health Check:

```
GET http://localhost:8080/health
```

The health check endpoint checks the status of:
- HTTP server
- MySQL database connection
- Redis connection

Response examples:

**All healthy:**
```json
{
  "status": "ok",
  "mysql": "healthy",
  "redis": "healthy"
}
```

**Degraded (one service down):**
```json
{
  "status": "degraded",
  "mysql": "healthy",
  "redis": "unhealthy"
}
```

7. API Documentation (Swagger):

```
GET http://localhost:8080/swagger/index.html
```

Access the interactive API documentation at `/swagger/index.html`

## Directory Structure
- `cmd/`: Main application entrypoint
- `internal/todo/`: Domain and service code for todo
  - `internal/todo/repository/`: Database repository layer (uses MySQL adapter)
- `internal/adapter/`: Infrastructure adapters
  - `internal/adapter/mysql`: MySQL database adapter (reusable across services)
  - `internal/adapter/redis`: Redis stream layer
- `internal/handler/http`: HTTP handlers
- `pkg/`: Reusable packages
  - `pkg/validator`: Request validation package
  - `pkg/errors`: Custom error types
  - `pkg/logger`: Structured logging with zap
- `docs/`: Swagger/OpenAPI documentation (generated)

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.

### Running Migrations

To run all pending migrations:

```sh
make migrate
```

Or directly:

```sh
go run ./cmd/main.go -migrate
```

### Migration Files

Migration files are located in `internal/migration/mysql/` and follow the naming convention:
- `001_create_todos.up.sql` - Creates the todos table
- `002_create_outbox.up.sql` - Creates the outbox table

### Notes

- Migrations are automatically applied when you run `make migrate`
- Make sure the database is running before executing migrations
- The migration tool will track which migrations have been applied in the database
- If a migration fails, you may need to manually fix the database state

## Development Mode

Run the application in development mode with colored logs:

```sh
make run-dev
```

Or:

```sh
go run ./cmd/main.go -dev
```

In development mode, logs are:
- Colored for better readability
- More verbose
- Human-readable format

## API Documentation

Generate Swagger documentation:

```sh
make swagger
```

Or:

```sh
swag init -g cmd/main.go -o docs
```

Then access the interactive API documentation at:
```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:
- Interactive API testing
- Request/response schemas
- Example requests
- Error response documentation

## Logging

The application uses structured logging with [zap](https://github.com/uber-go/zap):

- **Development mode**: Colored, human-readable logs
- **Production mode**: JSON formatted logs with timestamps
- **HTTP requests**: Logged with method, path, status, latency, and IP
- **Errors**: Structured error logging with context

Example log output:
```
INFO    HTTP request    {"method": "POST", "path": "/todo", "status": 201, "latency": "2.5ms", "ip": "127.0.0.1"}
```

## Testing

```sh
make test
```

## Benchmark

```sh
make benchmark
```

## Error Handling

The API uses structured error responses:

```json
{
  "code": 400,
  "message": "validation failed: description: is required"
}
```

Error codes:
- `400`: Bad Request (validation errors, invalid input)
- `404`: Not Found
- `500`: Internal Server Error

## Features

- ✅ Clean Architecture
- ✅ MySQL database with migrations
- ✅ Redis Stream integration
- ✅ Outbox pattern implementation
- ✅ Echo web framework
- ✅ Request validation
- ✅ UUID generation
- ✅ Graceful shutdown
- ✅ Health check endpoint with database/redis status
- ✅ Structured error handling
- ✅ Docker support with volumes for data persistence
- ✅ Structured logging with zap
- ✅ Swagger/OpenAPI documentation

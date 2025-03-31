# Gorilla/Mux + Graceful Shuwdown

This project is an HTTP API built using `mux`, following the Clean Architecture pattern with a PostgreSQL database.

## Features
- **Mux Routing**: Uses `mux` for defining API routes.
- **Clean Architecture**: Decoupled layers (Delivery, Usecase, Repository, Domain).
- **Graceful Shutdown**: Allows ongoing requests to complete before termination.
- **MySQL Database**: Stores vocabulary cards.

## Installation
### Clone the repository
```sh
git clone <repository-url>
cd <repository-name>
```

### Download dependencies
```sh
go mod tidy
```

### Download specific dependency
```sh
go get github.com/stretchr/testify
```

## Running the Server
```sh
go run main.go --port=8080
```

## Running Tests
```sh
go test -v
```

## API Endpoints
| Method | Endpoint     | Description         |
|--------|-------------|---------------------|
| GET    | `/cards`    | Retrieve all cards |

## How It Works
1. **Mux Routing**: The project uses `mux` to define API routes.
2. **Graceful Shutdown**: The server listens for termination signals and shuts down gracefully, allowing in-flight requests to complete.
3. **Logging**: Structured logging is enabled via `slog`.

## License
MIT

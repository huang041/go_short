# Go Short - URL Shortener Service

A simple and efficient URL shortener service built with Go, Gin, PostgreSQL, and Redis.

## Features

- Shorten long URLs to compact, easy-to-share links
- Multiple URL shortening algorithms (Base62, Base64, MD5, Random)
- Redis caching for improved performance
- RESTful API for URL management
- PostgreSQL database for persistent storage
- Docker and Docker Compose support for easy deployment

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM
- **Caching**: Redis
- **Architecture**: Domain-Driven Design (DDD) / Hexagonal Architecture
- **Database Migrations**: golang-migrate
- **Containerization**: Docker and Docker Compose
- **Configuration**: Environment variables via .env file

## Project Structure

```text
go_short/
├── .env                # Local environment variables (gitignored)
├── .env.example        # Example environment variables file
├── .gitignore          # Git ignore file
├── Dockerfile          # Production Docker image definition
├── Dockerfile.dev      # Development Docker image definition
├── README.md           # This file
├── architecture.md     # Mermaid diagram of the architecture
├── conf/               # Configuration management (loading .env)
├── domain/             # Core domain logic (business rules, entities)
│   ├── identity/       # User management domain
│   │   ├── entity/     # User entity definition
│   │   ├── repository/ # Repository interfaces (contracts for data access)
│   │   └── service/    # Identity domain services (e.g., user status logic)
│   └── urlshortener/   # URL shortener domain
│       ├── entity/     # URLMapping entity definition
│       ├── repository/ # Repository interfaces
│       └── service/    # URL shortening/lookup logic
├── go.mod              # Go module file
├── go.sum              # Go module checksum
├── infra/              # Infrastructure implementations (database, cache)
│   ├── database/       # Database connection setup (using GORM)
│   └── persistence/    # Repository implementations (gorm, redis)
├── internal/           # Internal application code (not meant for external reuse)
│   ├── api/            # API layer (routing, handlers)
│   │   ├── handler/    # HTTP request handlers (parsing, validation, calling app layer)
│   │   └── router.go   # API route definitions using Gin
│   ├── application/    # Application services (use cases orchestration)
│   │   ├── identity/   # Identity use cases (RegisterUser, AuthenticateUser)
│   │   └── urlshortener/ # URL shortener use cases
│   └── bootstrap/      # Dependency injection setup and application wiring
├── main.go             # Application entry point & lifecycle management
├── migrations/         # Database migration files (*.up.sql, *.down.sql)
└── scripts/            # Utility scripts (migration tool, dev shell)
```

## API Endpoints

### URL Shortener

-   `GET /ping` - Health check endpoint
-   `GET /url_mapping` - Get all URL mappings (consider adding pagination/filtering later)
-   `POST /url_mapping` - Create a new short URL (JSON body: `{"url": "...", "expires_in": <hours>}`)
-   `GET /{short_url}` - Redirect to the original URL

### User Authentication

-   `POST /auth/register` - Register a new user (JSON body: `{"username": "...", "email": "...", "password": "..."}`)
-   `POST /auth/login` - Log in a user (JSON body: `{"username": "...", "password": "..."}`)

## URL Shortening Algorithms

The service supports multiple URL shortening algorithms that can be configured via the `SHORTENER_ALGORITHM` environment variable in your `.env` file:

-   `base62` (default) - Converts database ID to a Base62 string (0-9, a-z, A-Z)
-   `base64` - Encodes the URL with Base64 and takes the first 8 characters
-   `md5` - Creates an MD5 hash of the URL + ID and takes the first 8 characters
-   `random` - Generates 8 random characters

## Getting Started

### Prerequisites

-   Docker and Docker Compose

### Development Environment (Recommended)

This setup uses Docker to provide a consistent development environment without needing Go, PostgreSQL, or Redis installed locally. Your local code is mounted into the container for live editing.

1.  **Clone the repository**
    ```bash
    git clone https://github.com/yourusername/go_short.git
    cd go_short
    ```
2.  **Create `.env` file**
    ```bash
    cp .env.example .env
    # Edit .env if needed, but defaults should work fine with Docker Compose network
    ```
3.  **Build and Start Development Containers**
    ```bash
    docker compose build dev # Only needed the first time or if Dockerfile.dev changes
    docker compose up -d dev postgres redis # Starts dev container, DB, and Cache
    ```
4.  **Run Database Migrations**
    ```bash
    ./scripts/migrate_tool.sh up
    ```
    *   This applies the latest database schema defined in the `migrations/` folder. Run this command whenever you add or change migration files.
5.  **Enter the Development Container**
    ```bash
    ./scripts/dev.sh
    ```
    *   This gives you an interactive shell (`sh`) inside the `go_short_dev` container. Your local project directory is mounted at `/app`.
6.  **Run the Application (inside the dev container)**
    *   **With hot-reload (using Air):**
        ```bash
        # (Inside dev container at /app)
        air # Air watches for file changes and automatically rebuilds/restarts
        ```
        The service will be available at `http://localhost:9081` (port mapped in `docker-compose.yml` for `dev`).
    *   **Manually:**
        ```bash
        # (Inside dev container at /app)
        go run main.go
        ```
7.  **Run Tests (inside the dev container)**
    ```bash
    # (Inside dev container at /app)
    go test ./...
    ```
8.  **Fix Go Modules (inside the dev container or using script)**
    If you add/remove dependencies, update `go.mod` and `go.sum`:
    ```bash
    # Option 1: Inside dev container
    go mod tidy

    # Option 2: From your host machine
    ./scripts/fix_go_mod.sh
    ```

### Production Deployment

1.  **Build and Start Production Containers**
    ```bash
    docker compose build app # Ensure the production image is built
    docker compose up -d app postgres redis # Start production app and dependencies
    ```
2.  **Run Database Migrations**
    This needs to be done against your production database *before* or *as* the new application version starts handling traffic.
    *   **Method 1: Using the migrate tool in the `app` container:** (The production `Dockerfile` includes the `migrate` tool and migration files)
        ```bash
        docker exec -it go_short_app sh -c "/usr/local/bin/migrate -path /app/migrations -database 'postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable' up"
        # Replace env vars with actual values or ensure they are available in the container
        ```
    *   **Method 2: Using the `migrate_tool.sh` script against the `dev` container:** (Connects to the same DB if running)
        ```bash
        # Ensure .env points to the correct database if not using the shared Docker network DB
        ./scripts/migrate_tool.sh up
        ```
    *   **Method 3: Integrate into your CI/CD pipeline:** Use the `migrate` CLI tool directly against the production DB connection string during deployment.
3.  The service will be available at `http://localhost:9080` (or the port mapped for the `app` service).

## Database Migrations

This project uses `golang-migrate` (version 4) to manage database schema changes in a versioned manner. SQL migration files are located in the `/migrations` directory. Each migration has an `up.sql` (apply change) and a `down.sql` (revert change) file.

Use the provided script (`./scripts/migrate_tool.sh`) to run migrations within the development container (it executes the `migrate` command inside the `dev` container):

-   **Apply all pending migrations:**
    ```bash
    ./scripts/migrate_tool.sh up
    ```
-   **Apply the next N migrations:**
    ```bash
    ./scripts/migrate_tool.sh up N
    ```
-   **Rollback the last migration:**
    ```bash
    ./scripts/migrate_tool.sh down 1
    ```
-   **Rollback all migrations:**
    ```bash
    ./scripts/migrate_tool.sh down
    ```
-   **Check the current migration version and dirty status:**
    ```bash
    ./scripts/migrate_tool.sh version
    ```
-   **Go to a specific version:**
    ```bash
    ./scripts/migrate_tool.sh goto <version_number>
    ```
-   **Force a specific version (marks previous migrations as clean/dirty, use with caution):**
    ```bash
    ./scripts/migrate_tool.sh force <version_number>
    ```

## Environment Variables

Configure the application using environment variables, typically defined in a `.env` file (see `.env.example`).

| Variable            | Description                      | Default    |
| ------------------- | -------------------------------- | ---------- |
| DB_HOST             | PostgreSQL host                  | postgres   |
| DB_PORT             | PostgreSQL port                  | 5432       |
| DB_USER             | PostgreSQL username              | postgres   |
| DB_PASSWORD         | PostgreSQL password              | postgres   |
| DB_NAME             | PostgreSQL database name         | go_short   |
| SHORTENER_ALGORITHM | URL shortening algorithm         | base62     |
| REDIS_HOST          | Redis host                       | redis      |
| REDIS_PORT          | Redis port                       | 6379       |
| REDIS_PASSWORD      | Redis password                   |            |
| REDIS_DB            | Redis database number            | 0          |
| GIN_MODE            | Gin framework mode (debug/release) | debug      |

## How It Works (High Level)

The application follows a layered architecture inspired by DDD and Hexagonal Architecture:

1.  **API Layer (`internal/api`)**: Receives HTTP requests (via Gin Router), validates basic input, and calls the appropriate Application Service.
2.  **Application Layer (`internal/application`)**: Orchestrates the steps needed to fulfill a use case (e.g., Register User, Create Short URL). It interacts with Domain Services and Repository Interfaces.
3.  **Domain Layer (`domain`)**: Contains the core business logic.
    *   **Entities**: Represent core concepts (User, URLMapping) with their state and behavior.
    *   **Services**: Encapsulate business logic that doesn't naturally fit within a single entity.
    *   **Repository Interfaces**: Define contracts for data persistence, abstracting away the database details.
4.  **Infrastructure Layer (`infra`)**: Provides concrete implementations for interfaces defined in other layers.
    *   **Persistence**: Implements Repository Interfaces using GORM (PostgreSQL) and go-redis.
    *   **Database**: Handles database connection setup.
    *   **Bootstrap**: Wires all the dependencies together on application startup.
5.  **Migrations (`migrations`)**: Manage database schema evolution using golang-migrate.

**Example Flow (Create Short URL):**

`HTTP POST /url_mapping` -> `API Handler` -> `URL Application Service` -> `URL Domain Service` (generates short code) -> `URL Repository Interface` -> `GORM Repository Implementation` -> `PostgreSQL`

**Example Flow (Redirect):**

`HTTP GET /{short_url}` -> `API Handler` -> `URL Application Service` -> `Cache Repository Interface` (check cache) -> (If cache miss) `URL Repository Interface` -> `GORM Repository Implementation` -> `PostgreSQL` -> `API Handler` (sends redirect)

## Performance Optimization

The service uses Redis caching primarily for the short URL to original URL lookup:

-   When a redirect request comes in, the system first checks Redis.
-   If the mapping is found in Redis (cache hit), the original URL is returned immediately, avoiding a database query.
-   If not found (cache miss), the system queries PostgreSQL, stores the result in Redis with a Time-To-Live (TTL, e.g., 24 hours), and then returns the original URL.
-   This significantly reduces database load for frequently accessed short URLs.
-   The system is designed to function even if Redis is temporarily unavailable (it will fall back to querying the database directly).

## License

This project is licensed under the MIT License - see the LICENSE file for details.
```
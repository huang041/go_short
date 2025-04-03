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
- **Containerization**: Docker and Docker Compose
- **Configuration**: Environment variables via .env file

## Project Structure

```
go_short/
├── conf/           # Configuration management
├── controllers/    # HTTP request handlers
├── models/         # Database models
├── routers/        # API route definitions
├── services/       # Business logic
├── .env.example    # Example environment variables
├── .gitignore      # Git ignore file
├── docker-compose.yml # Docker Compose configuration
├── Dockerfile      # Docker image definition
├── go.mod          # Go module file
├── go.sum          # Go module checksum
└── main.go         # Application entry point
```

## API Endpoints

- `GET /ping` - Health check endpoint
- `GET /url_mapping` - Get all URL mappings
- `POST /url_mapping` - Create a new short URL (JSON body with URL)
- `GET /{short_url}` - Redirect to the original URL

## URL Shortening Algorithms

The service supports multiple URL shortening algorithms that can be configured via the `SHORTENER_ALGORITHM` environment variable:

- `base62` (default) - Converts database ID to a Base62 string (0-9, a-z, A-Z)
- `base64` - Encodes the URL with Base64 and takes the first 8 characters
- `md5` - Creates an MD5 hash of the URL + ID and takes the first 8 characters
- `random` - Generates 8 random characters

## Getting Started

### Prerequisites

- Go 1.19 or higher
- PostgreSQL
- Redis
- Docker and Docker Compose (optional)

### Local Development

1. Clone the repository
   ```
   git clone https://github.com/yourusername/go_short.git
   cd go_short
   ```

2. Create and configure the .env file
   ```
   cp .env.example .env
   # Edit .env with your database and Redis credentials
   ```

3. Run the application
   ```
   go run main.go
   ```

4. The service will be available at `http://localhost:8080`

### Docker Deployment

1. Build and start the containers
   ```
   docker-compose up -d
   ```

2. The service will be available at `http://localhost:9080`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | PostgreSQL host | postgres |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL username | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | PostgreSQL database name | go_short |
| SHORTENER_ALGORITHM | URL shortening algorithm | base62 |
| REDIS_HOST | Redis host | redis |
| REDIS_PORT | Redis port | 6379 |
| REDIS_PASSWORD | Redis password | |
| REDIS_DB | Redis database number | 0 |

## How It Works

1. When a new URL is submitted, it's stored in the database
2. The system generates a short URL using the configured algorithm
3. The URL mapping is cached in Redis for faster access
4. When a short URL is accessed, the system:
   - First checks the Redis cache for the original URL
   - If not found in cache, queries the database
   - Redirects to the original URL
   - Caches the result for future requests (24-hour TTL)

## Performance Optimization

The service uses Redis caching to improve performance:
- Frequently accessed URLs are served from cache, reducing database load
- Cache entries expire after 24 hours to ensure data freshness
- The system gracefully degrades if Redis is unavailable (falls back to database)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

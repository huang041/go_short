# Go Short - URL Shortener Service

A simple and efficient URL shortener service built with Go, Gin, and PostgreSQL.

## Features

- Shorten long URLs to compact, easy-to-share links
- Base62 encoding for efficient URL shortening
- RESTful API for URL management
- PostgreSQL database for persistent storage
- Docker and Docker Compose support for easy deployment

## Tech Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM
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
- `POST /url_mapping?url={original_url}` - Create a new short URL

## Getting Started

### Prerequisites

- Go 1.16 or higher
- PostgreSQL
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
   # Edit .env with your database credentials
   ```

3. Run the application
   ```
   go run main.go
   ```

4. The service will be available at `http://127.0.0.1:8080`

### Docker Deployment

1. Build and start the containers
   ```
   docker-compose up -d
   ```

2. The service will be available at `http://localhost:8080`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | PostgreSQL host | postgres |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL username | postgres |
| DB_PASSWORD | PostgreSQL password | postgres |
| DB_NAME | PostgreSQL database name | go_short |

## How It Works

1. When a new URL is submitted, it's stored in the database
2. The database ID is converted to a Base62 string (using characters 0-9, a-z, A-Z)
3. This Base62 string becomes the short URL identifier
4. When a short URL is accessed, the system looks up the original URL and redirects

## License

This project is licensed under the MIT License - see the LICENSE file for details.

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

```
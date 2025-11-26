# Go Auth Service

A backend service providing Authentication with JWT tokens using Go (Fiber), PostgreSQL (GORM), and Redis.

## Overview

This project implements a secure authentication system following Clean Architecture principles. It supports user registration, login, token refresh, and logout with token blacklisting.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: Fiber
- **Database**: PostgreSQL
- **ORM**: GORM
- **Cache/Session**: Redis (Optional, used for token blacklisting)
- **Authentication**: JWT (Access + Refresh Tokens)

## Architecture

The project follows Clean Architecture with the following layers:

- **Domain (`internal/domain`)**: Core entities and interfaces (User, Repository interfaces).
- **Usecase (`internal/usecase`)**: Business logic (Auth logic).
- **Repository (`internal/repository`)**: Data access layer (GORM implementation).
- **Delivery (`internal/delivery/http`)**: HTTP handlers and routes (Fiber).
- **Infrastructure (`internal/infrastructure`)**: External adapters (DB, Redis).
- **Service (`internal/service`)**: Business logic (Password hashing,JWT).

## Prerequisites

- Go 1.20+
- Docker & Docker Compose

## Setup & Run

### Using Docker Compose (Recommended)

1. Clone the repository.
2. Create a `.env` file based on `.env.example`.
   ```bash
   cp .env.example .env
   ```
3. Run the services:
   ```bash
   docker compose up --build
   ```
   This will start PostgreSQL, Redis, and the Go API service.

### Running Locally

1. Ensure PostgreSQL and Redis are running.
2. Update `.env` with your local database credentials.
3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

### Authentication

- **Register**
  - `POST /auth/register`
  - Body: `{"email": "user@example.com", "password": "password", "name": "John Doe"}`

- **Login**
  - `POST /auth/login`
  - Body: `{"email": "user@example.com", "password": "password"}`
  - Returns: `access_token`, `refresh_token`

- **Refresh Token**
  - `POST /auth/refresh`
  - Body: `{"refresh_token": "..."}`
  - Returns: New `access_token`, New `refresh_token`

- **Logout**
  - `POST /auth/logout`
  - Headers: `Authorization: Bearer <access_token>`
  - Body: `{"refresh_token": "..."}`
  - Description: Invalidates both access and refresh tokens by blacklisting them in Redis.

### Protected Resources

- **Get Current User**
  - `GET /me`
  - Headers: `Authorization: Bearer <access_token>`
  - Returns: User profile information.

## Design Decisions

- **Clean Architecture**: Decouples business logic from frameworks and drivers, making the code testable and maintainable.
- **JWT**: Used for stateless authentication. Access tokens are short-lived (15m), refresh tokens are long-lived (24h).
- **Redis**: Used to store blacklisted tokens. This allows for immediate revocation of tokens upon logout, addressing a common JWT limitation.
- **GORM**: Used for database interactions to simplify SQL operations and migrations.
- **Fiber**: High-performance web framework for Go.

## Testing

To run tests (if added):
```bash
go test ./...
```

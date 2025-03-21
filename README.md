# Auth REST API Service

This is a simple authentication service built with Go that provides user registration, authentication with JWT tokens, and token management capabilities.

## Features

- User registration (sign up) with email and password
- User authentication (sign in) with JWT tokens
- JWT token verification and authorization
- Token revocation mechanism
- Token refresh mechanism
- In-memory storage for users and tokens

## Dependencies

- Go 1.21+
- Docker (for containerized deployment)

## Getting Started

### Running with Docker Compose (Recommended)

The easiest way to start the service is using Docker Compose:

```bash
docker-compose up
```

This will build the Docker image and start the service on port 8080.

### Running Locally

If you prefer to run the service locally:

1. Install Go 1.21 or higher
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build and run the service:
   ```bash
   go build -o auth-service
   ./auth-service
   ```

## API Endpoints

### 1. Sign Up

Register a new user with email and password.

```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### 2. Sign In

Authenticate a user and receive access and refresh tokens.

```bash
curl -X POST http://localhost:8080/api/auth/signin \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### 3. Verify Token

Verify that a token is valid (requires authentication).

```bash
curl -X GET http://localhost:8080/api/auth/verify \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 4. Get User Info

Get information about the authenticated user.

```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5. Refresh Token

Get a new access token using a refresh token.

```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_REFRESH_TOKEN"}'
```

### 6. Revoke Token

Revoke (invalidate) an access token.

```bash
curl -X POST http://localhost:8080/api/auth/revoke \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

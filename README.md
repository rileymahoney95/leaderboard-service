# Leaderboard Service

A RESTful API service for managing leaderboards with JWT-based authentication.

## Authentication

This service uses JSON Web Tokens (JWT) for authentication. The middleware implementation follows best practices:

- **Token Format**: Bearer token in the Authorization header
- **Token Validation**: Server-side validation of token signature and expiration
- **Role-Based Access Control**: Different endpoints require different user roles
- **Secure Token Storage**: Tokens should be stored securely on the client-side
- **Token Expiration**: Tokens have a configurable expiry time

### Authentication Flow

1. Client sends credentials to `/auth/login`
2. Server validates credentials and issues a JWT token
3. Client includes the token in the Authorization header for subsequent requests
4. Server validates the token for protected endpoints

### Example Login Request

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "role": "admin"
}
```

### Example Authenticated Request

```bash
curl -X GET http://localhost:8080/leaderboards \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## API Endpoints

### Public Endpoints

- `GET /health`: Health check endpoint
- `POST /auth/login`: Authenticate and get JWT token
- `POST /auth/register`: Register a new user (not implemented yet)

### Protected Endpoints (require authentication)

#### Available to all authenticated users

- `GET /leaderboards`: List all leaderboards
- `GET /leaderboards/{id}`: Get a specific leaderboard

#### Admin/Moderator only

- `POST /leaderboards`: Create a new leaderboard
- `PUT /leaderboards/{id}`: Update a leaderboard
- `DELETE /leaderboards/{id}`: Delete a leaderboard

## Environment Variables

Configure the following environment variables:

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/leaderboard_service?sslmode=disable
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRATION_HOURS=24
```

## Development

### Prerequisites

- Go 1.16+
- PostgreSQL

### Setup

1. Clone the repository
2. Create a `.env` file with the necessary environment variables
3. Run the application:

```bash
go run main.go
```

### Testing

```bash
go test ./...
```

## API Documentation

This service includes Swagger API documentation. After starting the server, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

### Authorization in Swagger UI

To test authenticated endpoints in Swagger UI:

1. First, use the `/auth/login` endpoint to get a JWT token
2. Click the "Authorize" button at the top of the page
3. In the value field, enter your token in the format: `Bearer YOUR_TOKEN_HERE`
4. Click "Authorize" and close the modal
5. Now you can access the authenticated endpoints

### Generating Swagger Documentation

If you make changes to the API, you need to regenerate the Swagger documentation:

```bash
swag init -g main.go
```

This command reads the annotations in your code and generates updated documentation.

## Troubleshooting

### Regenerating Swagger Documentation

If you encounter issues with Swagger documentation not displaying all routes or getting errors during swagger generation, follow these steps:

1. Install the Swagger CLI tool if not already installed:

   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. If you encounter errors with complex JSON objects in example annotations, simplify the examples. For metadata fields and other complex objects, use empty strings as examples instead of JSON objects:

   ```go
   // Instead of this:
   // example:"{\"country\":\"USA\",\"age\":25}"

   // Use this:
   // example:""
   ```

3. Run the Swagger initialization command:

   ```bash
   ~/go/bin/swag init
   ```

   Or if the swag command is in your PATH:

   ```bash
   swag init
   ```

4. Verify the generated files in the `docs` directory:

   - `docs/docs.go`
   - `docs/swagger.json`
   - `docs/swagger.yaml`

5. Restart your application and access the Swagger UI at:
   ```
   http://localhost:8080/swagger/index.html
   ```

These steps will ensure that all your routes, including new endpoints for participants, leaderboard entries, and metrics, are properly documented and accessible in the Swagger UI.

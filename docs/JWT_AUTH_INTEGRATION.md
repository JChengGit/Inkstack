# JWT Authentication Integration Guide

## Overview

This document describes the JWT authentication integration between the Auth Service and API Service in the Inkstack microservices architecture.

## Architecture

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│   Client    │         │ Auth Service│         │ API Service │
│             │         │  (port 8082)│         │ (port 8081) │
└──────┬──────┘         └──────┬──────┘         └──────┬──────┘
       │                       │                       │
       │  1. POST /register    │                       │
       │  or /login            │                       │
       │──────────────────────>│                       │
       │                       │                       │
       │  2. JWT tokens        │                       │
       │<──────────────────────│                       │
       │                       │                       │
       │  3. POST /posts       │                       │
       │  Authorization: Bearer JWT                    │
       │──────────────────────────────────────────────>│
       │                       │                       │
       │                       │   4. Validate JWT     │
       │                       │   (local signature)   │
       │                       │                       │
       │  5. Post created      │                       │
       │<──────────────────────────────────────────────│
```

## Services

### Auth Service (port 8082)
- **Database**: `auth_db` (PostgreSQL on port 5433)
- **Purpose**: User authentication, JWT token management
- **Endpoints**:
  - `POST /api/auth/register` - Create new user
  - `POST /api/auth/login` - Authenticate and get tokens
  - `POST /api/auth/refresh` - Refresh access token
  - `POST /api/auth/logout` - Revoke tokens
  - `GET /api/auth/me` - Get user profile
  - `POST /api/auth/change-password` - Change password
  - `POST /api/auth/validate` - Validate token (for API service)

### API Service (port 8081)
- **Database**: `api_db` (PostgreSQL on port 5432)
- **Purpose**: Business logic (posts, comments)
- **Authentication**: Validates JWT tokens from Auth Service
- **Endpoints**:
  - **Public** (no auth): GET /api/posts, GET /api/posts/:id
  - **Protected** (requires JWT): POST /api/posts, PUT /api/posts/:id, DELETE /api/posts/:id

### Redis (port 6379)
- **Purpose**: Token blacklist, rate limiting

## JWT Token Structure

### Access Token (15 minutes)
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "username": "testuser",
  "role": "user",
  "exp": 1234567890,
  "iat": 1234567890,
  "iss": "inkstack-auth"
}
```

### Refresh Token (7 days)
- Same structure as access token but longer expiry
- Stored in database for revocation tracking

## Security Features

### 1. Password Security
- Bcrypt hashing with cost factor 12
- Password strength requirements:
  - Minimum 8 characters
  - At least 1 uppercase letter
  - At least 1 lowercase letter
  - At least 1 number
  - At least 1 special character

### 2. Token Security
- JWT signed with HMAC-SHA256
- Shared secret (32+ characters) between services
- Access tokens: short-lived (15 min)
- Refresh tokens: long-lived (7 days), stored in DB, can be revoked

### 3. Rate Limiting
- Login attempts: 5 failures per 15 minutes
- Implemented via Redis counters

### 4. Token Blacklist
- Revoked tokens added to Redis blacklist
- TTL = remaining token lifetime
- Checked during validation

## Configuration

### Critical: JWT_SECRET Must Match!

Both services **must** use the identical `JWT_SECRET`:

**auth/.env**:
```env
JWT_SECRET=your-super-secret-jwt-key-change-in-production-min-32-chars
```

**api/.env**:
```env
JWT_SECRET=your-super-secret-jwt-key-change-in-production-min-32-chars
```

### Auth Service (.env)
```env
APP_ENV=dev
APP_PORT=8082

DB_HOST=localhost
DB_PORT=5433
DB_NAME=auth

JWT_SECRET=your-super-secret-jwt-key-change-in-production-min-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

REDIS_HOST=localhost
REDIS_PORT=6379
```

### API Service (.env)
```env
APP_ENV=dev
APP_PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_NAME=api

JWT_SECRET=your-super-secret-jwt-key-change-in-production-min-32-chars
AUTH_SERVICE_URL=http://localhost:8082
```

## Running the Services

### Option 1: Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

### Option 2: Local Development

**Terminal 1 - Auth Service:**
```bash
cd auth
go run cmd/server/main.go
```

**Terminal 2 - API Service:**
```bash
cd api
go run cmd/server/main.go
```

**Terminal 3 - PostgreSQL (Auth DB):**
```bash
docker run -d -p 5433:5432 \
  -e POSTGRES_DB=auth \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  postgres:15-alpine
```

**Terminal 4 - PostgreSQL (API DB):**
```bash
docker run -d -p 5432:5432 \
  -e POSTGRES_DB=api \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  postgres:15-alpine
```

**Terminal 5 - Redis:**
```bash
docker run -d -p 6379:6379 redis:7-alpine
```

## Testing the Integration

### Automated E2E Test

```bash
# Make script executable (Linux/Mac)
chmod +x test-e2e.sh

# Run test
./test-e2e.sh
```

### Manual Testing with cURL

#### 1. Register a User
```bash
curl -X POST http://localhost:8082/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "testuser",
    "password": "SecurePass123!@#"
  }'
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "testuser",
    "display_name": "",
    "role": "user"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### 2. Login (if already registered)
```bash
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "testuser",
    "password": "SecurePass123!@#"
  }'
```

#### 3. Create a Post (Authenticated)
```bash
curl -X POST http://localhost:8081/api/posts \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the post content",
    "excerpt": "A brief summary",
    "slug": "my-first-post"
  }'
```

**Response:**
```json
{
  "id": 1,
  "title": "My First Post",
  "slug": "my-first-post",
  "content": "This is the post content",
  "author_id": 1,
  "status": "draft",
  "created_at": "2026-01-03T15:30:00Z"
}
```

**Note**: `author_id` is automatically extracted from the JWT token!

#### 4. Try Without Authentication (Should Fail)
```bash
curl -X POST http://localhost:8081/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "This Should Fail",
    "content": "No token provided"
  }'
```

**Response:**
```json
{
  "error": "Authorization header required"
}
```

#### 5. Get User Profile
```bash
curl -X GET http://localhost:8082/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

#### 6. Refresh Token
```bash
curl -X POST http://localhost:8082/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

#### 7. Logout
```bash
curl -X POST http://localhost:8082/api/auth/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## API Documentation

### Swagger/OpenAPI

**Auth Service:**
- http://localhost:8082/swagger/index.html

**API Service:**
- http://localhost:8081/swagger/index.html

## Implementation Details

### How Auth Middleware Works

1. **Extract Token**: Parse `Authorization: Bearer <token>` header
2. **Validate JWT**: Verify signature using `JWT_SECRET`
3. **Extract Claims**: Get user_id, email, username, role from token
4. **Store in Context**: Save claims in Gin context for handlers
5. **Continue**: Pass request to handler

```go
// In API service middleware
func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        claims, err := jwtService.ValidateToken(token)

        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)

        c.Next()
    }
}
```

### How Handlers Use JWT Claims

```go
// In post handler
func (h *PostHandler) CreatePost(c *gin.Context) {
    // Extract user ID from context (set by middleware)
    userID, _ := c.Get("user_id")

    // Use userID as author_id
    post.AuthorID = userID.(uint)

    // Create post...
}
```

### Microservices Best Practices

#### ✅ What We Did Right

1. **Separate Databases**: Auth and API services have independent databases
2. **No Foreign Keys Across Services**: `author_id` and `user_id` are plain integers
3. **Local Token Validation**: API service validates JWT locally (fast, no network call)
4. **Shared Secret**: Both services use same `JWT_SECRET` for signing/validation
5. **Stateless Authentication**: JWT contains all needed info (user_id, role)

#### ⚠️ Trade-offs

1. **No Real-time Blacklist Check**: API service doesn't query Redis blacklist
   - **Pro**: Faster (no network call)
   - **Con**: Revoked tokens work until expiry
   - **Solution**: Keep access token TTL short (15 min)

2. **No Referential Integrity**: Database can't enforce user existence
   - **Pro**: Services are decoupled
   - **Con**: Can have posts by deleted users
   - **Solution**: Handle at application level

## Troubleshooting

### Problem: "JWT_SECRET is required"
**Solution**: Ensure `.env` file exists in both `api/` and `auth/` directories with matching `JWT_SECRET`

### Problem: "Invalid or expired token"
**Solution**:
- Check if access token expired (15 min TTL)
- Use refresh token to get new access token
- Verify `JWT_SECRET` matches in both services

### Problem: "User not authenticated" despite valid token
**Solution**:
- Ensure auth middleware is applied to the route
- Check if user_id is being extracted correctly in handler

### Problem: "Authorization header required"
**Solution**: Include header: `Authorization: Bearer YOUR_TOKEN`

### Problem: Docker containers can't connect
**Solution**:
- Ensure all services are on same Docker network
- Use service names (e.g., `auth-service:8082`) instead of `localhost` in Docker
- Check `depends_on` conditions are met

## Migration Notes

### Changes Made

1. **Removed** `api/internal/models/user.go` - Users now live in auth service
2. **Added** `api/internal/service/jwt_service.go` - JWT validation
3. **Added** `api/internal/middleware/auth_middleware.go` - Auth middleware
4. **Updated** `api/internal/config/config.go` - Added JWT config
5. **Updated** `api/cmd/server/main.go` - Applied middleware to routes
6. **Updated** `api/internal/handler/post_handler.go` - Extract user_id from JWT
7. **Updated** `api/internal/handler/comment_handler.go` - Extract user_id from JWT

### Database Schema

**Posts Table** (`api_db`):
- `author_id` (INTEGER) - References user in auth service (no FK constraint)

**Comments Table** (`api_db`):
- `user_id` (INTEGER) - References user in auth service (no FK constraint)

**Users Table** (`auth_db`):
- Lives in separate database
- Accessed only by auth service

## Production Checklist

- [ ] Change `JWT_SECRET` to strong random value (32+ chars)
- [ ] Set `JWT_SECRET` identically in both services
- [ ] Use HTTPS in production
- [ ] Set `APP_ENV=prod` in both services
- [ ] Use secure database passwords
- [ ] Enable Redis password authentication
- [ ] Set up proper logging and monitoring
- [ ] Configure CORS for production origins
- [ ] Use secrets manager (e.g., AWS Secrets Manager) for sensitive values
- [ ] Set up database backups
- [ ] Configure health check endpoints
- [ ] Use proper SSL certificates

## Future Enhancements

1. **OAuth 2.0**: Google, GitHub login
2. **Email Verification**: Verify email addresses
3. **2FA**: Two-factor authentication
4. **Password Reset**: Email-based password reset
5. **Account Lockout**: Lock after multiple failed attempts
6. **Audit Logging**: Track authentication events
7. **API Keys**: Service-to-service authentication
8. **Redis Blacklist in API**: Check revoked tokens in real-time

## Support

- **Swagger Docs**: http://localhost:8082/swagger/ and http://localhost:8081/swagger/
- **Health Checks**: http://localhost:8082/health and http://localhost:8081/health
- **GitHub Issues**: Report bugs and feature requests

---

**Last Updated**: January 3, 2026
**Version**: 1.0.0

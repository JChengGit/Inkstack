# Auth Service - Microservice Setup

## GOAL
Create a standalone authentication microservice with user management, JWT token generation, and OAuth 2.0 support. This service will be completely separate from the main API service with its own database.

## ARCHITECTURE OVERVIEW

### Service Communication
```
Client → API Service (requires JWT) → validates token
         ↓
Client → Auth Service → issues JWT token
```

### Databases
- **auth_db** (port 5433): Users, passwords, refresh tokens, OAuth tokens
- **api_db** (port 5432): Posts, comments, business data
- **Redis** (port 6379): Session cache, token blacklist

### Ports
- Auth Service: 8082
- API Service: 8081

## PROJECT STRUCTURE

```
auth/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── migrate.go
│   ├── models/
│   │   ├── base.go
│   │   ├── user.go
│   │   ├── refresh_token.go
│   │   └── oauth_account.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── token_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── jwt_service.go
│   │   └── oauth_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── oauth_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── rate_limit.go
│   └── util/
│       ├── password.go
│       ├── validator.go
│       └── response.go
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_refresh_tokens.up.sql
│   └── 000002_create_refresh_tokens.down.sql
├── .env.example
├── .env
├── Dockerfile
├── go.mod
└── go.sum
```

## REQUIREMENTS

### 0. Migrate API Database - Remove Users Table

**IMPORTANT:** Before setting up the auth service, we need to remove the users table from the API database since users will now live in the auth database.

#### Step 0.1: Create Migration to Drop Users Table from API Database

Create `api/migrations/000005_drop_users_table.up.sql`:
```sql
-- Drop indexes first
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop users table
DROP TABLE IF EXISTS users;
```

Create `api/migrations/000005_drop_users_table.down.sql`:
```sql
-- Recreate users table (for rollback)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
```

#### Step 0.2: Update API Database Schema

Since posts and comments reference users by `author_id` and `user_id`, we need to handle these foreign keys:

**Option A: Remove Foreign Key Constraints (Recommended for microservices)**
- Posts and comments will store `author_id`/`user_id` as plain integers
- No database-level foreign key constraints
- User validation happens at application level via auth service

Create `api/migrations/000006_remove_user_foreign_keys.up.sql`:
```sql
-- Remove foreign key constraints from posts (if they exist)
ALTER TABLE IF EXISTS posts DROP CONSTRAINT IF EXISTS posts_author_id_fkey;

-- Remove foreign key constraints from comments (if they exist)
ALTER TABLE IF EXISTS comments DROP CONSTRAINT IF EXISTS comments_user_id_fkey;

-- author_id and user_id columns remain, just without FK constraints
```

Create `api/migrations/000006_remove_user_foreign_keys.down.sql`:
```sql
-- Recreate foreign keys (won't work in microservices architecture)
-- This is just for documentation
ALTER TABLE IF EXISTS posts
    ADD CONSTRAINT posts_author_id_fkey
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE IF EXISTS comments
    ADD CONSTRAINT comments_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
```

#### Step 0.3: Remove User Model from API Service

Delete or move the following file:
- `api/internal/models/user.go` - Move to auth service

**Note:** Keep the `author_id` and `user_id` fields in Post and Comment models. These will be populated from JWT tokens.

#### Step 0.4: Update API .env File

Update `api/.env`:
```env
APP_ENV=dev
APP_PORT=8081

# API Database (separate from auth)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inkstack_api
DB_SSLMODE=disable

# Auth Service Configuration
AUTH_SERVICE_URL=http://localhost:8082
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

**CRITICAL:** The `JWT_SECRET` must be the same in both API and Auth services!

### 1. Initialize Auth Service Project

Create new Go module in `auth/` directory:
```bash
cd auth
go mod init inkstack-auth
```

Add dependencies:
```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/redis/go-redis/v9
go get golang-migrate/migrate/v4
```

### 2. Environment Configuration

Create `.env.example`:
```env
# Application
APP_ENV=dev
APP_PORT=8082

# Database
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inkstack_auth
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# OAuth (for future)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
```

### 3. Database Models

#### User Model (`internal/models/user.go`)
```go
type User struct {
    BaseModel
    Email           string     `gorm:"uniqueIndex;not null" json:"email"`
    Username        string     `gorm:"uniqueIndex;not null" json:"username"`
    PasswordHash    string     `gorm:"not null" json:"-"` // Never expose in JSON
    DisplayName     string     `json:"display_name"`
    Bio             string     `gorm:"type:text" json:"bio"`
    AvatarURL       string     `json:"avatar_url"`
    EmailVerified   bool       `gorm:"default:false" json:"email_verified"`
    IsActive        bool       `gorm:"default:true" json:"is_active"`
    Role            string     `gorm:"default:'user'" json:"role"` // user, admin
    LastLoginAt     *time.Time `json:"last_login_at"`
}
```

#### RefreshToken Model (`internal/models/refresh_token.go`)
```go
type RefreshToken struct {
    BaseModel
    UserID    uint      `gorm:"not null;index" json:"user_id"`
    Token     string    `gorm:"uniqueIndex;not null" json:"token"`
    ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
    IsRevoked bool      `gorm:"default:false" json:"is_revoked"`
    IPAddress string    `json:"ip_address"`
    UserAgent string    `json:"user_agent"`
}
```

### 4. Database Migrations

#### Migration 000001: Users Table
```sql
-- 000001_create_users.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    bio TEXT,
    avatar_url VARCHAR(500),
    email_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    role VARCHAR(20) DEFAULT 'user',
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

#### Migration 000002: Refresh Tokens Table
```sql
-- 000002_create_refresh_tokens.up.sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token VARCHAR(500) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
```

### 5. JWT Service

Create `internal/service/jwt_service.go`:

**Token Structure:**
```go
type JWTClaims struct {
    UserID   uint   `json:"user_id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

**Methods:**
- `GenerateAccessToken(user *models.User) (string, error)` - 15min expiry
- `GenerateRefreshToken(user *models.User) (string, error)` - 7 days expiry
- `ValidateToken(tokenString string) (*JWTClaims, error)`
- `RefreshAccessToken(refreshToken string) (string, error)`
- `RevokeToken(token string) error` - Add to Redis blacklist

### 6. Auth Service

Create `internal/service/auth_service.go`:

**Methods:**
- `Register(email, username, password string) (*models.User, error)`
  - Validate email format
  - Check uniqueness
  - Hash password with bcrypt
  - Create user in database
  - Send verification email (future)

- `Login(emailOrUsername, password string) (accessToken, refreshToken string, error)`
  - Find user by email or username
  - Verify password
  - Check if account is active
  - Generate JWT tokens
  - Store refresh token
  - Update last_login_at

- `RefreshToken(refreshToken string) (newAccessToken string, error)`
  - Validate refresh token
  - Check if revoked
  - Check expiry
  - Generate new access token

- `Logout(refreshToken string) error`
  - Revoke refresh token
  - Add access token to blacklist

- `ChangePassword(userID uint, oldPassword, newPassword string) error`

### 7. HTTP Handlers

Create `internal/handler/auth_handler.go`:

**Endpoints:**

1. `POST /api/auth/register`
   - Request: `{ "email", "username", "password" }`
   - Response: `{ "user": {...}, "access_token": "...", "refresh_token": "..." }`

2. `POST /api/auth/login`
   - Request: `{ "email_or_username", "password" }`
   - Response: `{ "user": {...}, "access_token": "...", "refresh_token": "..." }`

3. `POST /api/auth/refresh`
   - Request: `{ "refresh_token": "..." }`
   - Response: `{ "access_token": "..." }`

4. `POST /api/auth/logout`
   - Header: `Authorization: Bearer <token>`
   - Request: `{ "refresh_token": "..." }`
   - Response: `{ "message": "Logged out successfully" }`

5. `POST /api/auth/change-password`
   - Header: `Authorization: Bearer <token>`
   - Request: `{ "old_password", "new_password" }`
   - Response: `{ "message": "Password changed successfully" }`

6. `GET /api/auth/me` - Get current user profile
   - Header: `Authorization: Bearer <token>`
   - Response: `{ "user": {...} }`

7. `POST /api/auth/validate` - Validate token (for API service)
   - Request: `{ "token": "..." }`
   - Response: `{ "valid": true, "user_id": 123, "role": "user" }`

### 8. User Handler

Create `internal/handler/user_handler.go`:

**Endpoints:**
- `GET /api/users/:id` - Get user profile (public info only)
- `PUT /api/users/:id` - Update user profile (auth required)
- `GET /api/users` - List users (admin only)

### 9. Middleware

Create `internal/middleware/auth_middleware.go`:

**Middleware Functions:**
- `RequireAuth()` - Validates JWT token
- `RequireRole(role string)` - Check user role
- `RateLimit()` - Limit login attempts

### 10. Password Utilities

Create `internal/util/password.go`:

**Functions:**
- `HashPassword(password string) (string, error)` - bcrypt hash
- `ComparePassword(hashedPassword, password string) bool`
- `ValidatePasswordStrength(password string) error`
  - Min 8 characters
  - At least 1 uppercase
  - At least 1 lowercase
  - At least 1 number
  - At least 1 special character

### 11. Redis Integration

Create `internal/database/redis.go`:

**Usage:**
- Token blacklist (logout)
- Rate limiting (login attempts)
- Session cache (optional)

**Functions:**
- `BlacklistToken(token string, expiry time.Duration) error`
- `IsTokenBlacklisted(token string) (bool, error)`
- `IncrementLoginAttempts(email string) (int, error)`
- `ResetLoginAttempts(email string) error`

### 12. API Service Integration

Update API service to validate tokens with auth service:

Create `api/internal/middleware/auth_middleware.go`:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Validate token locally using JWT_SECRET
        // OR call auth service to validate
        token := strings.TrimPrefix(authHeader, "Bearer ")

        // Option 1: Validate locally (faster)
        claims, err := jwtService.ValidateToken(token)

        // Option 2: Call auth service (more secure, can check blacklist)
        // valid, userID := callAuthService(token)

        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Store user info in context
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)

        c.Next()
    }
}
```

Apply to protected routes:
```go
posts := api.Group("/posts")
posts.Use(middleware.AuthMiddleware())
{
    posts.POST("", postHandler.CreatePost) // Now requires auth
}
```

### 13. Dockerfile for Auth Service

Create `auth/Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ../auth/todo .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env* ./

EXPOSE 8082

CMD ["./main"]
```

## TESTING CHECKLIST

### Registration & Login
- [ ] Register new user with valid data
- [ ] Register with duplicate email (should fail)
- [ ] Register with duplicate username (should fail)
- [ ] Register with weak password (should fail)
- [ ] Login with email and password
- [ ] Login with username and password
- [ ] Login with wrong password (should fail)
- [ ] Login with non-existent user (should fail)
- [ ] Receive access and refresh tokens on successful login

### Token Management
- [ ] Access token expires after 15 minutes
- [ ] Refresh token works to get new access token
- [ ] Expired refresh token is rejected
- [ ] Revoked token cannot be used
- [ ] Logout invalidates refresh token
- [ ] Token validation endpoint works

### Password Management
- [ ] Change password with correct old password
- [ ] Change password fails with wrong old password
- [ ] Password is hashed in database (never plaintext)

### Protected Routes
- [ ] Access protected endpoint with valid token
- [ ] Access protected endpoint without token (401)
- [ ] Access protected endpoint with expired token (401)
- [ ] Access protected endpoint with invalid token (401)

### Integration with API Service
- [ ] API service validates tokens from auth service
- [ ] API service rejects invalid tokens
- [ ] User ID from token is used in API operations
- [ ] Posts/comments are associated with authenticated user

### Rate Limiting
- [ ] Multiple failed login attempts are rate limited
- [ ] Successful login resets rate limit counter

## DELIVERABLES

1. Complete auth service project structure
2. User registration and login functionality
3. JWT token generation and validation
4. Refresh token mechanism
5. Password hashing and validation
6. Database migrations for users and tokens
7. Redis integration for blacklist
8. Auth middleware for API service
9. Docker setup for both services
10. Updated docker-compose.yml with separate databases

## SECURITY CONSIDERATIONS

1. **Password Security**
   - Use bcrypt with cost factor 12+
   - Never log passwords
   - Never return password hashes in API responses

2. **JWT Security**
   - Use strong secret (min 32 characters)
   - Set appropriate expiry times
   - Store refresh tokens securely
   - Implement token rotation

3. **API Security**
   - HTTPS only in production
   - Rate limiting on auth endpoints
   - CORS configuration
   - Input validation and sanitization

4. **Database Security**
   - Separate databases for auth and API
   - Use prepared statements (GORM handles this)
   - Regular backups
   - Encrypt sensitive data at rest

## COMMUNICATION PATTERNS

### Pattern 1: Token Validation in API Service (Recommended)
API service validates JWT locally using shared secret:
- **Pros:** Fast, no network call
- **Cons:** Cannot check token blacklist in real-time

### Pattern 2: Auth Service Validation
API service calls auth service to validate token:
- **Pros:** Can check blacklist, more secure
- **Cons:** Slower, network dependency

### Pattern 3: Hybrid
Validate JWT locally + check Redis blacklist in API service:
- **Pros:** Fast and secure
- **Cons:** API service needs Redis access

**Recommendation:** Use Pattern 3 for production.

## FUTURE ENHANCEMENTS

- OAuth 2.0 integration (Google, GitHub, etc.)
- Email verification workflow
- Two-factor authentication (2FA)
- Password reset via email
- Account lockout after failed attempts
- Audit logging for security events
- API key management for service-to-service auth

## NOTES

- Keep auth service simple and focused
- Auth service should be stateless (except for database/redis)
- Use Redis for ephemeral data (blacklist, rate limits)
- Consider using a secrets manager in production
- Document all API endpoints with OpenAPI/Swagger
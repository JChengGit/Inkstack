# Inkstack Microservices Architecture

## Overview

Inkstack is now architected as a microservices system with two main services:

```
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Client    │────────▶│ Auth Service │────────▶│  Auth DB    │
│             │         │   :8082      │         │   :5433     │
└─────────────┘         └──────────────┘         └─────────────┘
       │                        │
       │                   JWT Token
       │                        │
       ▼                        ▼
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│  Validates  │────────▶│ API Service  │────────▶│  API DB     │
│  JWT Token  │         │   :8081      │         │   :5432     │
└─────────────┘         └──────────────┘         └─────────────┘
                               │
                               ▼
                        ┌─────────────┐
                        │   Redis     │
                        │   :6379     │
                        └─────────────┘
```

## Services

### 1. Auth Service (Port 8082)
**Responsibility:** User authentication and authorization
- User registration & login
- JWT token generation & validation
- Password management
- OAuth 2.0 integration (future)
- Token refresh & revocation

**Database:** `inkstack_auth` (Port 5433)
- Users table
- Refresh tokens
- OAuth accounts (future)

**Tech Stack:**
- Go + Gin
- PostgreSQL
- Redis (token blacklist, rate limiting)
- JWT (github.com/golang-jwt/jwt/v5)
- Bcrypt (password hashing)

### 2. API Service (Port 8081)
**Responsibility:** Business logic and data management
- Posts CRUD
- Comments CRUD
- Tags, search, media (future)

**Database:** `inkstack_api` (Port 5432)
- Posts table
- Comments table
- Tags, media (future)

**Tech Stack:**
- Go + Gin
- PostgreSQL
- GORM

### 3. Redis (Port 6379)
**Responsibility:** Caching and temporary data
- JWT token blacklist
- Rate limiting counters
- Session cache

## Quick Start

### 1. Start All Services with Docker Compose

```bash
cd api/local/Inkstack_containers
docker-compose up -d
```

This starts:
- `api_db` (PostgreSQL :5432)
- `auth_db` (PostgreSQL :5433)
- `redis` (Redis :6379)

### 2. Development Without Docker

#### Terminal 1 - Start Databases
```bash
cd api/local
docker-compose up api_db auth_db redis
```

#### Terminal 2 - Run Auth Service
```bash
cd auth
cp .env.example .env
# Edit .env with your configuration
go run cmd/server/main.go
```

#### Terminal 3 - Run API Service
```bash
cd api
cp .env.example .env
# Edit .env with your configuration
go run cmd/server/main.go
```

## API Flow

### User Registration
```bash
# 1. Register with Auth Service
POST http://localhost:8082/api/auth/register
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "SecurePass123!"
}

# Response:
{
  "user": { "id": 1, "email": "user@example.com", ... },
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci..."
}
```

### User Login
```bash
# 2. Login
POST http://localhost:8082/api/auth/login
{
  "email_or_username": "johndoe",
  "password": "SecurePass123!"
}

# Response:
{
  "user": { ... },
  "access_token": "eyJhbGci...",
  "refresh_token": "eyJhbGci..."
}
```

### Create Post (Protected)
```bash
# 3. Create post with JWT token
POST http://localhost:8081/api/posts
Headers:
  Authorization: Bearer eyJhbGci...
Body:
{
  "title": "My First Post",
  "content": "This is the content..."
}

# API service validates token and creates post
# User ID is extracted from JWT token
```

### Token Refresh
```bash
# 4. When access token expires (15 min)
POST http://localhost:8082/api/auth/refresh
{
  "refresh_token": "eyJhbGci..."
}

# Response:
{
  "access_token": "eyJhbGci..."  # New access token
}
```

## Database Separation

### Why Separate Databases?

1. **Security Isolation**
   - Auth data (passwords, tokens) is isolated
   - Compromise of API DB doesn't expose auth data

2. **Scalability**
   - Scale each database independently
   - Auth DB: Read-heavy (token validation)
   - API DB: Read/Write mixed (posts, comments)

3. **Service Independence**
   - Each service owns its data
   - No direct cross-database queries
   - Clean microservice boundaries

4. **Performance**
   - Different optimization strategies
   - Separate connection pools
   - Independent backup/restore

### Database Access Rules

✅ **Allowed:**
- Auth service → Auth DB
- API service → API DB
- Each service reads its own database

❌ **Not Allowed:**
- API service → Auth DB (use API calls instead)
- Auth service → API DB
- Direct cross-database joins

### Data Synchronization

When API service needs user info:

**Option 1: Embed in JWT (Recommended)**
```json
{
  "user_id": 123,
  "username": "johndoe",
  "email": "user@example.com",
  "role": "user"
}
```
User info is in the token, no need to query auth service.

**Option 2: Call Auth Service**
```bash
GET http://auth-service:8082/api/users/123
```
When you need fresh user data (e.g., check if user is still active).

**Option 3: Denormalization**
Store minimal user info in API DB (username, avatar_url) for display purposes. Sync via events (future).

## Environment Variables

### Auth Service (.env)
```env
APP_PORT=8082
DB_HOST=localhost
DB_PORT=5433
DB_NAME=inkstack_auth
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
REDIS_HOST=localhost
REDIS_PORT=6379
```

### API Service (.env)
```env
APP_PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_NAME=inkstack_api
JWT_SECRET=your-secret-key  # Same as auth service!
AUTH_SERVICE_URL=http://localhost:8082
```

**Important:** `JWT_SECRET` must be the same in both services for token validation.

## Security Considerations

### JWT Token Flow
1. User logs in → Auth service generates JWT
2. Client includes JWT in requests to API service
3. API service validates JWT using shared secret
4. API service extracts user_id from token
5. No need to call auth service on every request

### Token Structure
```json
{
  "user_id": 123,
  "email": "user@example.com",
  "username": "johndoe",
  "role": "user",
  "exp": 1234567890,  // Expiry timestamp
  "iat": 1234567890   // Issued at
}
```

### Token Types
- **Access Token:** Short-lived (15 min), used for API requests
- **Refresh Token:** Long-lived (7 days), used to get new access tokens

### Token Revocation
When user logs out:
1. Client sends refresh token to auth service
2. Auth service marks token as revoked in database
3. Auth service adds token to Redis blacklist
4. API service checks blacklist before accepting tokens

## Implementation Steps

1. ✅ **Setup Docker Compose** (Done)
   - Separate databases
   - Redis instance
   - Network configuration

2. 🔲 **Implement Auth Service** (See `auth/todo/1-auth-service-init.md`)
   - User registration & login
   - JWT token generation
   - Token refresh & revocation
   - Password management

3. 🔲 **Update API Service**
   - Add JWT validation middleware
   - Extract user_id from tokens
   - Protect routes with authentication
   - Remove user management from API service

4. 🔲 **Testing**
   - Test registration & login
   - Test token validation
   - Test protected routes
   - Test token refresh & revocation

5. 🔲 **Documentation**
   - API documentation (Swagger)
   - Deployment guide
   - Security best practices

## Migration Strategy

If you have existing users in API service:

1. **Create auth service with empty database**
2. **Migration script:**
   ```sql
   -- Export users from api_db
   COPY (SELECT id, email, username, created_at FROM users) TO '/tmp/users.csv' CSV HEADER;

   -- Import to auth_db (set temporary passwords)
   COPY users FROM '/tmp/users.csv' CSV HEADER;
   ```
3. **Send password reset emails** to all migrated users
4. **Keep user references** in api_db for foreign keys

## Troubleshooting

### Port Conflicts
```bash
# Check what's using the port
lsof -i :5432  # or :5433, :8081, :8082
```

### Database Connection Issues
```bash
# Check if containers are running
docker ps

# Check logs
docker logs inkstack_auth_db
docker logs inkstack_api_db
```

### Token Validation Failures
- Ensure `JWT_SECRET` is the same in both services
- Check token hasn't expired
- Verify token format: `Authorization: Bearer <token>`

## Next Steps

1. Read `auth/todo/1-auth-service-init.md` for detailed implementation guide
2. Start by implementing basic registration & login
3. Add JWT token generation
4. Integrate with API service
5. Add OAuth 2.0 support (future)

## References

- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OAuth 2.0 RFC](https://tools.ietf.org/html/rfc6749)
- [Microservices Patterns](https://microservices.io/patterns/index.html)
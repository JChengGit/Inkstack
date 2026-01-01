# Swagger Quick Start Guide

## Quick Access

### Auth Service Swagger UI
```
http://localhost:8082/swagger/index.html
```

### API Service Swagger UI
```
http://localhost:8081/swagger/index.html
```

## Common Tasks

### 1. Start Services and View Documentation

```bash
# Terminal 1 - Start Auth Service
cd auth
go run cmd/server/main.go

# Terminal 2 - Start API Service
cd api
go run cmd/server/main.go

# Open in browser:
# Auth: http://localhost:8082/swagger/index.html
# API:  http://localhost:8081/swagger/index.html
```

### 2. Test Authentication Flow

#### Step 1: Register a User
```bash
curl -X POST http://localhost:8082/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "SecurePass123!"
  }'
```

#### Step 2: Login
```bash
curl -X POST http://localhost:8082/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email_or_username": "test@example.com",
    "password": "SecurePass123!"
  }'
```

Copy the `access_token` from the response.

#### Step 3: Access Protected Endpoint
```bash
curl -X GET http://localhost:8082/api/auth/me \
  -H "Authorization: Bearer <your-access-token>"
```

### 3. Test API Endpoints

#### Create a Post
```bash
curl -X POST http://localhost:8081/api/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first post.",
    "excerpt": "A brief summary",
    "author_id": 1
  }'
```

#### Get All Posts
```bash
curl http://localhost:8081/api/posts
```

#### Create a Comment
```bash
curl -X POST http://localhost:8081/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Great post!"
  }'
```

### 4. Regenerate Documentation After Changes

```bash
# For Auth Service
cd auth
swag init -g cmd/server/main.go --output docs

# For API Service
cd api
swag init -g cmd/server/main.go --output docs
```

## Swagger Annotation Cheat Sheet

### Basic Endpoint
```go
// @Summary Brief description
// @Description Detailed description
// @Tags category-name
// @Accept json
// @Produce json
// @Success 200 {object} ResponseType
// @Router /api/endpoint [get]
func Handler(c *gin.Context) {}
```

### With Path Parameter
```go
// @Param id path int true "Resource ID"
// @Router /api/posts/{id} [get]
```

### With Query Parameter
```go
// @Param page query int false "Page number" default(1)
// @Param status query string false "Filter by status"
// @Router /api/posts [get]
```

### With Request Body
```go
// @Param request body CreateRequest true "Request body"
// @Router /api/posts [post]
```

### With Authentication
```go
// @Security BearerAuth
// @Router /api/auth/me [get]
```

### Multiple Response Codes
```go
// @Success 200 {object} ResponseType
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/posts/{id} [get]
```

## Testing in Swagger UI

### 1. Using the UI
1. Click on an endpoint to expand it
2. Click "Try it out"
3. Fill in parameters
4. Click "Execute"
5. View response

### 2. With Authentication
1. Get token from login/register
2. Click "Authorize" button (top right)
3. Enter: `Bearer <token>`
4. Click "Authorize"
5. Test protected endpoints

## HTTP Status Codes Used

| Code | Meaning | Usage |
|------|---------|-------|
| 200 | OK | Successful GET, PUT |
| 201 | Created | Successful POST |
| 204 | No Content | Successful DELETE |
| 400 | Bad Request | Validation errors |
| 401 | Unauthorized | Missing/invalid auth |
| 404 | Not Found | Resource not found |
| 500 | Internal Error | Server errors |
| 503 | Service Unavailable | Health check failed |

## Common Issues

### Port Already in Use
```bash
# Find process using port
# Windows:
netstat -ano | findstr :8081
taskkill /PID <PID> /F

# Linux/Mac:
lsof -i :8081
kill -9 <PID>
```

### Swagger UI Not Loading
1. Check service is running: `curl http://localhost:8081/health`
2. Ensure docs generated: `ls auth/docs/` or `ls api/docs/`
3. Verify import in main.go: `_ "service-name/docs"`

### Changes Not Showing
1. Regenerate docs: `swag init -g cmd/server/main.go --output docs`
2. Restart service
3. Hard refresh browser (Ctrl+F5)

## Response Format Examples

### Success Response
```json
{
  "message": "Success",
  "data": {
    "id": 1,
    "title": "Example"
  }
}
```

### Error Response
```json
{
  "error": "Validation failed",
  "message": "Email is required"
}
```

### Paginated Response
```json
{
  "posts": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_pages": 5
  }
}
```

## Request Examples

### Auth Service

**Register:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "SecurePass123!"
}
```

**Login:**
```json
{
  "email_or_username": "user@example.com",
  "password": "SecurePass123!"
}
```

**Change Password:**
```json
{
  "old_password": "OldPass123!",
  "new_password": "NewPass456!"
}
```

### API Service

**Create Post:**
```json
{
  "title": "My Blog Post",
  "content": "Full content here...",
  "excerpt": "Brief summary",
  "slug": "my-blog-post",
  "author_id": 1
}
```

**Update Post:**
```json
{
  "title": "Updated Title",
  "status": "published"
}
```

**Create Comment:**
```json
{
  "content": "This is a comment",
  "parent_id": 5
}
```

## Environment Setup

Make sure these environment variables are set:

### Auth Service (.env)
```env
APP_PORT=8082
DB_HOST=localhost
DB_PORT=5433
DB_NAME=inkstack_auth
JWT_SECRET=your-secret-key
```

### API Service (.env)
```env
APP_PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_NAME=inkstack_api
```

## Next Steps

1. ✅ View API documentation in Swagger UI
2. ✅ Test authentication flow
3. ✅ Create posts and comments
4. ⬜ Integrate auth middleware in API service
5. ⬜ Add more endpoints as needed
6. ⬜ Deploy to production

## Resources

- Full documentation: `/swag/README.md`
- Swagger UI: Browse to service URL + `/swagger/index.html`
- OpenAPI specs: `auth/docs/swagger.json` and `api/docs/swagger.json`
- Handler source: `auth/internal/handler/*` and `api/internal/handler/*`

# Inkstack Swagger/OpenAPI Documentation

This directory contains documentation and information about the Swagger/OpenAPI implementation for Inkstack microservices.

## Overview

Inkstack uses [swaggo/swag](https://github.com/swaggo/swag) to automatically generate OpenAPI 3.0 documentation from code annotations. Both the **Auth Service** and **API Service** have fully documented REST APIs with interactive Swagger UI.

## Architecture

```
┌─────────────────────┐         ┌─────────────────────┐
│   Auth Service      │         │    API Service      │
│   Port: 8082        │         │    Port: 8081       │
│                     │         │                     │
│  /swagger/*         │         │  /swagger/*         │
│  (Swagger UI)       │         │  (Swagger UI)       │
└─────────────────────┘         └─────────────────────┘
```

## Accessing Swagger UI

### Auth Service
- **URL**: `http://localhost:8082/swagger/index.html`
- **Description**: Authentication and user management endpoints
- **Key Features**:
  - User registration and login
  - JWT token management
  - Password operations
  - Token validation (for API service)

### API Service
- **URL**: `http://localhost:8081/swagger/index.html`
- **Description**: Blog posts and comments API
- **Key Features**:
  - Posts CRUD operations
  - Comments CRUD operations
  - Post publishing workflow
  - Comment moderation

## Generated Files

Each service has its documentation generated in the `docs/` directory:

```
auth/docs/
├── docs.go         # Go code for embedding swagger
├── swagger.json    # OpenAPI JSON specification
└── swagger.yaml    # OpenAPI YAML specification

api/docs/
├── docs.go         # Go code for embedding swagger
├── swagger.json    # OpenAPI JSON specification
└── swagger.yaml    # OpenAPI YAML specification
```

## Implementation Details

### Dependencies

Both services use the following packages:
- `github.com/swaggo/swag` - Swagger generator
- `github.com/swaggo/gin-swagger` - Gin middleware for Swagger UI
- `github.com/swaggo/files` - Embedded Swagger UI files

### Annotation Format

#### General API Information (in `main.go`)
```go
// @title Inkstack Auth Service API
// @version 1.0
// @description Authentication and user management microservice
// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

#### Endpoint Documentation (in handlers)
```go
// @Summary Register a new user
// @Description Create a new user account with email, username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // handler implementation
}
```

### Security

Auth Service implements Bearer token authentication:
- Protected endpoints require `Authorization: Bearer <token>` header
- Swagger UI includes authorization input for testing protected endpoints
- Token validation endpoint is available for service-to-service communication

## Regenerating Documentation

When you make changes to API handlers or annotations, regenerate the documentation:

### For Auth Service:
```bash
cd auth
swag init -g cmd/server/main.go --output docs
```

### For API Service:
```bash
cd api
swag init -g cmd/server/main.go --output docs
```

Or use the full path if `swag` is not in your PATH:
```bash
$(go env GOPATH)/bin/swag init -g cmd/server/main.go --output docs
```

## API Endpoints

### Auth Service (`localhost:8082`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | No |
| POST | `/api/auth/login` | User login | No |
| POST | `/api/auth/refresh` | Refresh access token | No |
| POST | `/api/auth/validate` | Validate JWT token | No |
| GET | `/api/auth/me` | Get current user profile | Yes |
| POST | `/api/auth/logout` | Logout user | Yes |
| POST | `/api/auth/change-password` | Change password | Yes |
| GET | `/health` | Health check | No |

### API Service (`localhost:8081`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| **Posts** | | | |
| GET | `/api/posts` | List posts (paginated) | No |
| POST | `/api/posts` | Create new post | No* |
| GET | `/api/posts/{id}` | Get post by ID | No |
| GET | `/api/posts/slug/{slug}` | Get post by slug | No |
| PUT | `/api/posts/{id}` | Update post | No* |
| DELETE | `/api/posts/{id}` | Delete post | No* |
| POST | `/api/posts/{id}/publish` | Publish post | No* |
| POST | `/api/posts/{id}/unpublish` | Unpublish post | No* |
| **Comments** | | | |
| GET | `/api/posts/{id}/comments` | List comments for post | No |
| POST | `/api/posts/{id}/comments` | Create comment on post | No* |
| GET | `/api/comments/{id}` | Get comment by ID | No |
| PUT | `/api/comments/{id}` | Update comment | No* |
| DELETE | `/api/comments/{id}` | Delete comment | No* |
| POST | `/api/comments/{id}/approve` | Approve comment | No* |
| POST | `/api/comments/{id}/reject` | Reject comment | No* |
| **Health** | | | |
| GET | `/hello` | Hello world | No |
| GET | `/health` | Health check | No |

\* *Note: These endpoints will require authentication once auth middleware is fully integrated*

## Testing with Swagger UI

### Testing Protected Endpoints (Auth Service)

1. First, register or login to get an access token:
   - Use `/api/auth/login` or `/api/auth/register`
   - Copy the `access_token` from the response

2. Click the "Authorize" button in Swagger UI
3. Enter: `Bearer <your-access-token>`
4. Click "Authorize" then "Close"
5. Now you can test protected endpoints

### Testing API Endpoints

1. Open the Swagger UI for the respective service
2. Expand the endpoint you want to test
3. Click "Try it out"
4. Fill in required parameters
5. Click "Execute"
6. View the response

## Best Practices

### When Adding New Endpoints

1. **Add Swagger annotations** above the handler function:
   ```go
   // @Summary Short description
   // @Description Detailed description
   // @Tags category
   // @Accept json
   // @Produce json
   // @Param name location type required "description"
   // @Success 200 {object} ResponseType
   // @Failure 400 {object} map[string]interface{}
   // @Router /path [method]
   ```

2. **Document request/response models** with JSON tags:
   ```go
   type CreateUserRequest struct {
       Email    string `json:"email" binding:"required"`
       Username string `json:"username" binding:"required"`
   }
   ```

3. **Regenerate documentation**:
   ```bash
   swag init -g cmd/server/main.go --output docs
   ```

4. **Test in Swagger UI** to ensure it appears correctly

### Annotation Guidelines

- Use clear, concise summaries
- Provide detailed descriptions when necessary
- Group related endpoints with `@Tags`
- Document all parameters (path, query, body)
- Include all possible response codes
- Use proper types for request/response bodies

## Troubleshooting

### Swagger UI shows 404 or not loading
- Ensure the service is running
- Check that docs are generated in `docs/` directory
- Verify import: `_ "service-name/docs"` in main.go

### Changes not reflected in Swagger UI
- Regenerate documentation with `swag init`
- Restart the service
- Clear browser cache

### "swag: command not found"
```bash
go install github.com/swaggo/swag/cmd/swag@latest
# or use full path
$(go env GOPATH)/bin/swag init -g cmd/server/main.go --output docs
```

### Import cycle or build errors
- Ensure `docs` package is only imported in main.go with blank identifier `_`
- Don't import `docs` package in other files

## OpenAPI Specification Files

The generated OpenAPI specifications can be used with:
- **Postman**: Import `swagger.json` for collection generation
- **API Clients**: Generate SDKs using OpenAPI generators
- **Documentation Sites**: Host static documentation
- **Testing Tools**: Automated API testing

## Future Enhancements

- [ ] Add request/response examples in annotations
- [ ] Document rate limiting behavior
- [ ] Add OAuth 2.0 flow documentation (when implemented)
- [ ] Create unified API gateway documentation
- [ ] Add API versioning documentation
- [ ] Generate client SDKs automatically

## References

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Gin Swagger Integration](https://github.com/swaggo/gin-swagger)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

## Support

For issues or questions about the API documentation:
- Check the Swagger UI for endpoint details
- Review this README for setup instructions
- See individual handler files for implementation details
- Refer to TODO files in `/todo` directory for feature status

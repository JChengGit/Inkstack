# Suggested Updates for Main README.md

This file contains suggested updates to add to the main `/README.md` to document the Swagger implementation.

## Add to README.md

### Under the "Services" Section

After the existing service descriptions, add:

```markdown
### API Documentation

Both services include interactive Swagger/OpenAPI documentation:

- **Auth Service API Docs**: http://localhost:8082/swagger/index.html
- **API Service API Docs**: http://localhost:8081/swagger/index.html

The API documentation provides:
- Interactive endpoint testing
- Request/response schemas
- Authentication testing
- OpenAPI 3.0 specifications

For more details, see [Swagger Documentation](/swag/README.md).
```

### Add New Section: "API Documentation"

Add this new section after the Architecture section:

```markdown
## API Documentation

Inkstack uses Swagger/OpenAPI for API documentation. Each microservice has its own interactive documentation:

### Accessing Documentation

1. **Start the services**:
   ```bash
   # Terminal 1 - Auth Service
   cd auth && go run cmd/server/main.go

   # Terminal 2 - API Service
   cd api && go run cmd/server/main.go
   ```

2. **Open Swagger UI**:
   - Auth Service: http://localhost:8082/swagger/index.html
   - API Service: http://localhost:8081/swagger/index.html

### Quick Test

```bash
# Register a user
curl -X POST http://localhost:8082/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"Test123!"}'

# Create a post
curl -X POST http://localhost:8081/api/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello World","content":"My first post","author_id":1}'
```

### Documentation Resources

- 📖 [Complete Swagger Guide](/swag/README.md)
- 🚀 [Quick Start Guide](/swag/QUICK_START.md)
- 📋 [Implementation Summary](/swag/IMPLEMENTATION_SUMMARY.md)
```

### Update the Technology Stack Section

Under each service's tech stack, add:

```markdown
**Documentation:**
- Swagger/OpenAPI 3.0
- Gin-Swagger middleware
- Interactive API testing
```

### Add to the "Development" Section

```markdown
### API Documentation

Generate Swagger documentation after making changes to handlers:

```bash
# For Auth Service
cd auth
swag init -g cmd/server/main.go --output docs

# For API Service
cd api
swag init -g cmd/server/main.go --output docs
```

View documentation at:
- Auth: http://localhost:8082/swagger/index.html
- API: http://localhost:8081/swagger/index.html
```

## Alternative: Minimal Addition

If you prefer a minimal update, just add this to the README:

```markdown
## API Documentation

Interactive Swagger documentation is available for both services:
- Auth Service: http://localhost:8082/swagger/index.html
- API Service: http://localhost:8081/swagger/index.html

See [/swag/README.md](/swag/README.md) for complete documentation.
```

## Implementation

You can either:
1. Manually copy these sections into `/README.md`
2. Or ask me to update the main README with these additions

The choice depends on whether you want to maintain the existing README structure or integrate these changes more deeply into the existing sections.

# Database Connection Setup

## GOAL
Set up PostgreSQL database connection with proper configuration management, environment handling, ORM integration, health checks, and migration support.

## REQUIREMENTS

### 1. Dependencies
Add the following Go modules:
- `gorm.io/gorm` - ORM framework
- `gorm.io/driver/postgres` - PostgreSQL driver for GORM
- `github.com/joho/godotenv` - Environment variable loader
- `golang-migrate/migrate/v4` - Database migration tool
- `golang-migrate/migrate/v4/database/postgres` - PostgreSQL driver for migrate
- `golang-migrate/migrate/v4/source/file` - File source for migrations

### 2. Environment Configuration

#### Create `.env.example` (template)
```
# Application
APP_ENV=dev
APP_PORT=8081

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inkstack_dev
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

#### Create `.env` (local development)
Same structure as `.env.example` but with actual local credentials. This file should be gitignored.

#### Create `.env.prod.example` (production template)
```
APP_ENV=prod
APP_PORT=8080

DB_HOST=prod-db-host
DB_PORT=5432
DB_USER=inkstack_prod
DB_PASSWORD=secure_password_here
DB_NAME=inkstack_prod
DB_SSLMODE=require
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=10m
```

### 3. Configuration Package

Create `internal/config/config.go`:
- Define configuration structs for App and Database settings
- Load environment variables using godotenv
- Provide environment detection (development, staging, production)
- Validate required configuration values
- Support fallback to default values where appropriate

### 4. Database Package

Create `internal/database/postgres.go`:
- Initialize GORM connection with PostgreSQL
- Configure connection pool (max open/idle connections, lifetime)
- Implement connection retry logic with exponential backoff
- Provide database instance accessor (singleton or dependency injection)
- Implement graceful shutdown/close function
- Add proper error handling and logging

Create `internal/database/health.go`:
- Implement health check function that pings the database
- Return connection status and error details
- Include connection pool statistics (open/idle connections)

### 5. Database Models

Create `internal/models/base.go`:
- Define base model with common fields (ID, CreatedAt, UpdatedAt, DeletedAt)
- Use GORM conventions (soft delete support)

Create example model `internal/models/user.go`:
- Sample User model to verify ORM setup
- Include basic fields (ID, Email, Username, CreatedAt, UpdatedAt, DeletedAt)

Create demo table model `internal/models/test_table.go`:
- Model name: `TestTable`
- Fields:
  - `ID` (uint, primary key, auto-increment) - inherited from base model
  - `Foo` (string) - text field
  - `Bar` (int) - numeric field
  - `CreatedAt`, `UpdatedAt`, `DeletedAt` (timestamps) - inherited from base model
- Use GORM tags for proper column mapping
- Embed the base model for common fields

### 6. Database Migration

Create migration structure:
```
migrations/
├── 000001_init_schema.up.sql
└── 000001_init_schema.down.sql
```

Initial migration should:
- Create users table with proper constraints
- Add indexes for commonly queried fields
- Include timestamps and soft delete support

Create additional migration for test table:
```
migrations/
├── 000002_create_test_table.up.sql
└── 000002_create_test_table.down.sql
```

Test table migration should:
- Create `test_tables` table with columns:
  - `id` (SERIAL PRIMARY KEY)
  - `foo` (VARCHAR(255) or TEXT) - string field
  - `bar` (INTEGER) - numeric field
  - `created_at` (TIMESTAMP)
  - `updated_at` (TIMESTAMP)
  - `deleted_at` (TIMESTAMP, nullable) - for soft delete support
- Add any necessary indexes
- Down migration should drop the table

Create `internal/database/migrate.go`:
- Function to run migrations programmatically
- Support both up and down migrations
- Handle migration versioning
- Provide migration status checking

### 7. Main Application Integration

Update `cmd/api/main.go` or equivalent:
- Load configuration on startup
- Initialize database connection
- Run migrations (optional: flag to skip in production)
- Register database health check endpoint
- Ensure graceful shutdown closes database connection

### 8. Health Check Endpoint

Create `/health` or `/api/v1/health` endpoint:
- Returns JSON with application status
- Includes database connection status
- Returns appropriate HTTP status codes (200 for healthy, 503 for unhealthy)
- Example response:
```json
{
  "status": "healthy",
  "timestamp": "2025-10-22T10:30:00Z",
  "database": {
    "status": "connected",
    "open_connections": 3,
    "idle_connections": 2
  }
}
```

### 9. Error Handling & Logging

- Use structured logging (consider `go.uber.org/zap` or `github.com/sirupsen/logrus`)
- Log database connection events (success, failure, retry attempts)
- Log migration execution
- Provide meaningful error messages for configuration issues

### 10. Update .gitignore

Ensure the following are ignored:
```
.env
.env.local
.env.*.local
*.db
```

## TESTING CHECKLIST

- [ ] Application starts successfully with valid .env configuration
- [ ] Application fails gracefully with clear error message when database is unreachable
- [ ] Environment switching works (development vs production settings)
- [ ] Health check endpoint returns correct status
- [ ] Migrations run successfully on fresh database
- [ ] Connection pool settings are respected
- [ ] Graceful shutdown closes database connections properly
- [ ] Can perform basic CRUD operation with sample User model
- [ ] TestTable model is properly created and can be used for CRUD operations
- [ ] Test table migration (000002) runs successfully
- [ ] Can insert, read, update, and delete records from test_tables

## DELIVERABLES

1. Configuration management with environment detection
2. PostgreSQL connection with GORM
3. Connection pooling and retry logic
4. Health check endpoint
5. Migration system with initial schema
6. Sample model for verification (User model)
7. Demo TestTable model with migrations
8. Proper error handling and logging
9. Updated .gitignore
10. Documentation in .env.example files

## NOTES

- Follow Go project layout best practices
- Keep database logic in `internal/database` package
- Use dependency injection or singleton pattern consistently
- Ensure all database credentials are loaded from environment variables
- Never commit actual .env files with sensitive data
- Consider adding database query logging for development environment

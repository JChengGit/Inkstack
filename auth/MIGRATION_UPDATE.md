# Auth Service Migration Update

## Summary
Updated the Auth service to use **golang-migrate** with SQL migration files, matching the API service's implementation for automatic database migrations.

## What Changed

### 1. Updated `internal/database/migrate.go`

**Before:**
- Simple implementation with commented-out GORM AutoMigrate
- No parameters - `RunMigrations()`
- No actual migration execution

**After:**
- Full golang-migrate implementation matching API service
- Takes migrations path parameter - `RunMigrations(migrationsPath string)`
- Automatically runs SQL migration files
- Additional utility functions:
  - `RollbackMigration(migrationsPath string)` - Rollback last migration
  - `GetMigrationVersion(migrationsPath string)` - Get current version
  - `MigrateTo(migrationsPath string, version uint)` - Migrate to specific version
  - `toFileURL(path string)` - Cross-platform file URL conversion (Windows/Unix)

### 2. Updated `cmd/server/main.go`

**Before:**
```go
// Run migrations (optional - you can also use SQL migrations manually)
if err := database.RunMigrations(); err != nil {
    log.Printf("Warning: Failed to run migrations: %v", err)
}
```

**After:**
```go
// Run migrations
// Use relative path - requires running from auth/ directory
// Alternative: Use environment variable MIGRATIONS_PATH for flexibility
migrationsPath := "./migrations"

log.Printf("Running database migrations from: %s", migrationsPath)
if err := database.RunMigrations(migrationsPath); err != nil {
    log.Printf("Warning: Failed to run migrations: %v", err)
    log.Println("Continuing anyway - migrations may need to be run manually")
    log.Println("Note: Ensure you're running from the auth/ directory")
}
```

### 3. Added Dependencies

Automatically installed via `go mod tidy`:
- `github.com/golang-migrate/migrate/v4`
- `github.com/golang-migrate/migrate/v4/database/postgres`
- `github.com/golang-migrate/migrate/v4/source/file`

## Existing Migration Files

The auth service already has SQL migration files in place:

```
auth/migrations/
├── 000001_create_users.up.sql       # Creates users table
├── 000001_create_users.down.sql     # Drops users table
├── 000002_create_refresh_tokens.up.sql   # Creates refresh_tokens table
└── 000002_create_refresh_tokens.down.sql # Drops refresh_tokens table
```

These migrations will now be **automatically applied** when the auth service starts.

## How It Works

1. **Service Startup**: When the auth service starts, it:
   - Connects to PostgreSQL database
   - Connects to Redis
   - **Runs pending migrations automatically**
   - Continues with service initialization

2. **Migration Tracking**:
   - golang-migrate creates a `schema_migrations` table in the database
   - Tracks which migrations have been applied
   - Only runs new migrations that haven't been applied yet

3. **First Run**: On first run, it will apply:
   - `000001_create_users.up.sql` - Creates users table with indexes
   - `000002_create_refresh_tokens.up.sql` - Creates refresh_tokens table

4. **Subsequent Runs**:
   - Checks which migrations are already applied
   - Only runs new migrations (if any)
   - Logs "No new migrations to apply" if database is up to date

## Benefits

✅ **Consistent with API Service** - Both services now use the same migration approach
✅ **Production-Ready** - Uses industry-standard golang-migrate library
✅ **Automatic** - Migrations run automatically on service startup
✅ **Rollback Support** - Can rollback migrations if needed
✅ **Version Control** - SQL files tracked in Git
✅ **Cross-Platform** - Works on Windows, Linux, and macOS
✅ **Safe** - Only applies new migrations, never re-runs existing ones

## Testing

### Start Auth Service
```bash
cd auth
go run cmd/server/main.go
```

### Expected Output
```
2025/11/12 22:00:00 Starting Inkstack Auth Service in dev mode
2025/11/12 22:00:00 Connected to PostgreSQL: inkstack_auth
2025/11/12 22:00:00 Connected to Redis at localhost:6379
2025/11/12 22:00:00 Running database migrations from: ./migrations
2025/11/12 22:00:00 Migrations applied successfully
2025/11/12 22:00:00 Auth service starting on port 8082
```

Or if migrations are already applied:
```
2025/11/12 22:00:00 Running database migrations from: ./migrations
2025/11/12 22:00:00 No new migrations to apply
```

## Manual Migration Operations

### Rollback Last Migration
```go
import "inkstack-auth/internal/database"

err := database.RollbackMigration("./migrations")
```

### Get Current Version
```go
version, dirty, err := database.GetMigrationVersion("./migrations")
fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)
```

### Migrate to Specific Version
```go
err := database.MigrateTo("./migrations", 1) // Migrate to version 1
```

## Database Schema

After migrations run, the auth database will have:

### Users Table
- `id` - Primary key
- `email` - Unique, indexed
- `username` - Unique, indexed
- `password_hash` - Bcrypt hashed password
- `display_name` - Optional display name
- `bio` - User biography
- `avatar_url` - Profile picture URL
- `email_verified` - Email verification status
- `is_active` - Account active status
- `role` - User role (user, admin)
- `last_login_at` - Last login timestamp
- `created_at`, `updated_at`, `deleted_at` - Timestamps

### Refresh Tokens Table
- `id` - Primary key
- `user_id` - Foreign key to users (indexed)
- `token` - Unique token string (indexed)
- `expires_at` - Token expiration (indexed)
- `is_revoked` - Revocation status
- `ip_address` - Client IP address
- `user_agent` - Client user agent
- `created_at`, `updated_at`, `deleted_at` - Timestamps

### Schema Migrations Table (created automatically)
- `version` - Migration version number
- `dirty` - Migration state flag

## Adding New Migrations

To add new migrations in the future:

1. Create new migration files:
   ```bash
   # In auth/migrations/ directory
   # Next version number is 000003
   touch 000003_add_column_to_users.up.sql
   touch 000003_add_column_to_users.down.sql
   ```

2. Write SQL for up migration:
   ```sql
   -- 000003_add_column_to_users.up.sql
   ALTER TABLE users ADD COLUMN phone_number VARCHAR(20);
   CREATE INDEX idx_users_phone ON users(phone_number);
   ```

3. Write SQL for down migration:
   ```sql
   -- 000003_add_column_to_users.down.sql
   DROP INDEX IF EXISTS idx_users_phone;
   ALTER TABLE users DROP COLUMN phone_number;
   ```

4. Restart service - migration will apply automatically

## Comparison: Auth vs API Service

Both services now have identical migration setups:

| Feature | Auth Service | API Service |
|---------|--------------|-------------|
| Migration Library | golang-migrate ✅ | golang-migrate ✅ |
| Auto-run on startup | Yes ✅ | Yes ✅ |
| SQL migration files | Yes ✅ | Yes ✅ |
| Rollback support | Yes ✅ | Yes ✅ |
| Version tracking | Yes ✅ | Yes ✅ |
| Cross-platform | Yes ✅ | Yes ✅ |
| Migrations path | `./migrations` | `./migrations` |
| Implementation | Identical | Identical |

## Troubleshooting

### Migration Fails
If migration fails, check:
1. Database is running and accessible
2. Running from correct directory (`auth/`)
3. Migration files have correct permissions
4. SQL syntax is correct

### "No such file or directory" Error
- Ensure you're running from the `auth/` directory
- Or use absolute path for migrations

### Dirty Migration State
If a migration partially fails:
```bash
# Check version and dirty state
version, dirty, _ := database.GetMigrationVersion("./migrations")

# Fix the migration SQL
# Then force version (use with caution)
# Or manually fix schema_migrations table in database
```

## Notes

- **Production**: Same setup works in production
- **Docker**: Works in Docker containers (migrations run on container start)
- **Multiple Instances**: Safe to run multiple instances - golang-migrate uses database locks
- **Zero Downtime**: For zero-downtime deployments, consider forward-compatible migrations

## Related Files

- `/auth/internal/database/migrate.go` - Migration implementation
- `/auth/cmd/server/main.go` - Migration execution
- `/auth/migrations/*.sql` - SQL migration files
- `/api/internal/database/migrate.go` - API service (identical implementation)

---

**Status**: ✅ Complete and tested
**Updated**: 2025-11-12
**Auth Service**: Now matches API service migration implementation

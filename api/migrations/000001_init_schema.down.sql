-- Drop indexes
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop users table
DROP TABLE IF EXISTS users;
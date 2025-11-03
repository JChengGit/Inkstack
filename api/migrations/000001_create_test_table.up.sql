-- Create test_tables table
CREATE TABLE IF NOT EXISTS test_tables (
    id SERIAL PRIMARY KEY,
    foo VARCHAR(255),
    bar INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_test_tables_deleted_at ON test_tables(deleted_at);

-- Add comments
COMMENT ON TABLE test_tables IS 'Demo table for testing database connectivity and CRUD operations';
COMMENT ON COLUMN test_tables.foo IS 'Text field for testing';
COMMENT ON COLUMN test_tables.bar IS 'Numeric field for testing';
COMMENT ON COLUMN test_tables.deleted_at IS 'Soft delete timestamp';
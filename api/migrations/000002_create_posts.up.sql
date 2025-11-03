-- Create posts table
-- Note: author_id references users in the separate auth service database
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    author_id INTEGER NOT NULL,  -- References auth_db.users.id (no FK constraint in microservices)
    status VARCHAR(20) NOT NULL DEFAULT 'draft',  -- draft, published, archived
    published_at TIMESTAMP,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at);
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);

-- Add table and column comments
COMMENT ON TABLE posts IS 'Blog posts content';
COMMENT ON COLUMN posts.author_id IS 'User ID from auth service (no FK constraint - microservices architecture)';
COMMENT ON COLUMN posts.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN posts.status IS 'Post status: draft, published, archived';
COMMENT ON COLUMN posts.excerpt IS 'Short summary for post listing';

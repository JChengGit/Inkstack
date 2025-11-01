-- Create comments table
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    parent_id INTEGER,  -- For nested comments/replies
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending, approved, rejected, spam
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);

-- Add table and column comments
COMMENT ON TABLE comments IS 'User comments on posts';
COMMENT ON COLUMN comments.parent_id IS 'Parent comment ID for nested replies';
COMMENT ON COLUMN comments.status IS 'Moderation status: pending, approved, rejected, spam';

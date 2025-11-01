-- Drop indexes
DROP INDEX IF EXISTS idx_posts_deleted_at;
DROP INDEX IF EXISTS idx_posts_slug;
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_author_id;

-- Drop posts table
DROP TABLE IF EXISTS posts;

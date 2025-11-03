-- Drop indexes
DROP INDEX IF EXISTS idx_comments_deleted_at;
DROP INDEX IF EXISTS idx_comments_status;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;

-- Drop comments table
DROP TABLE IF EXISTS comments;

# Posts & Comments CRUD API

## GOAL
Implement posts and comments functionality with proper database schema, ORM models, repository pattern, service layer, and RESTful CRUD endpoints.

## REQUIREMENTS

### 1. Database Schema Design

#### Posts Table
Create migration: `migrations/000003_create_posts.up.sql`

```sql
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    author_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',  -- draft, published, archived
    published_at TIMESTAMP,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_posts_author_id ON posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at);
CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_at ON posts(deleted_at);

-- Comments
COMMENT ON TABLE posts IS 'Blog posts content';
COMMENT ON COLUMN posts.slug IS 'URL-friendly identifier';
COMMENT ON COLUMN posts.status IS 'Post status: draft, published, archived';
COMMENT ON COLUMN posts.excerpt IS 'Short summary for post listing';
```

Down migration: `migrations/000003_create_posts.down.sql`
```sql
DROP INDEX IF EXISTS idx_posts_deleted_at;
DROP INDEX IF EXISTS idx_posts_slug;
DROP INDEX IF EXISTS idx_posts_published_at;
DROP INDEX IF EXISTS idx_posts_status;
DROP INDEX IF EXISTS idx_posts_author_id;
DROP TABLE IF EXISTS posts;
```

#### Comments Table
Create migration: `migrations/000004_create_comments.up.sql`

```sql
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

-- Indexes
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comments_status ON comments(status);
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);

-- Comments
COMMENT ON TABLE comments IS 'User comments on posts';
COMMENT ON COLUMN comments.parent_id IS 'Parent comment ID for nested replies';
COMMENT ON COLUMN comments.status IS 'Moderation status: pending, approved, rejected, spam';
```

Down migration: `migrations/000004_create_comments.down.sql`
```sql
DROP INDEX IF EXISTS idx_comments_deleted_at;
DROP INDEX IF EXISTS idx_comments_status;
DROP INDEX IF EXISTS idx_comments_parent_id;
DROP INDEX IF EXISTS idx_comments_user_id;
DROP INDEX IF EXISTS idx_comments_post_id;
DROP TABLE IF EXISTS comments;
```

### 2. ORM Models

#### Post Model
Create `internal/models/post.go`:
- Embed BaseModel (ID, CreatedAt, UpdatedAt, DeletedAt)
- Fields:
  - Title (string, required, max 255)
  - Slug (string, unique, required, max 255)
  - Content (string, required)
  - Excerpt (string, optional)
  - AuthorID (uint, required, foreign key to users)
  - Status (string, enum: draft/published/archived)
  - PublishedAt (*time.Time, nullable)
  - ViewCount (int, default 0)
- GORM tags for proper mapping
- JSON tags for API responses
- Custom table name: "posts"
- Add validation tags

#### Comment Model
Create `internal/models/comment.go`:
- Embed BaseModel
- Fields:
  - PostID (uint, required, foreign key)
  - UserID (uint, required, foreign key)
  - ParentID (*uint, nullable, foreign key for replies)
  - Content (string, required)
  - Status (string, enum: pending/approved/rejected/spam)
- GORM tags with foreign key relationships
- JSON tags
- Custom table name: "comments"
- Add validation tags

### 3. Repository Layer

#### Post Repository
Create `internal/repository/post_repository.go`:

Interface methods:
```go
type PostRepository interface {
    Create(post *models.Post) error
    FindByID(id uint) (*models.Post, error)
    FindBySlug(slug string) (*models.Post, error)
    FindAll(limit, offset int) ([]models.Post, error)
    FindByAuthor(authorID uint, limit, offset int) ([]models.Post, error)
    FindByStatus(status string, limit, offset int) ([]models.Post, error)
    Update(post *models.Post) error
    Delete(id uint) error  // Soft delete
    IncrementViewCount(id uint) error
    Count() (int64, error)
    CountByAuthor(authorID uint) (int64, error)
}
```

Implementation:
- Use GORM for all database operations
- Handle errors properly
- Use proper GORM queries (Preload, Where, etc.)

#### Comment Repository
Create `internal/repository/comment_repository.go`:

Interface methods:
```go
type CommentRepository interface {
    Create(comment *models.Comment) error
    FindByID(id uint) (*models.Comment, error)
    FindByPostID(postID uint) ([]models.Comment, error)
    FindByUserID(userID uint, limit, offset int) ([]models.Comment, error)
    FindReplies(parentID uint) ([]models.Comment, error)
    Update(comment *models.Comment) error
    Delete(id uint) error  // Soft delete
    UpdateStatus(id uint, status string) error
    CountByPost(postID uint) (int64, error)
}
```

Implementation:
- Use GORM
- Handle cascade relationships
- Proper error handling

### 4. Service Layer

#### Post Service
Create `internal/service/post_service.go`:

Business logic methods:
```go
type PostService interface {
    CreatePost(title, content, excerpt, slug string, authorID uint) (*models.Post, error)
    GetPost(id uint) (*models.Post, error)
    GetPostBySlug(slug string) (*models.Post, error)
    ListPosts(page, pageSize int) ([]models.Post, int64, error)
    ListPostsByAuthor(authorID uint, page, pageSize int) ([]models.Post, int64, error)
    ListPublishedPosts(page, pageSize int) ([]models.Post, int64, error)
    UpdatePost(id uint, updates map[string]interface{}) (*models.Post, error)
    DeletePost(id uint) error
    PublishPost(id uint) (*models.Post, error)
    UnpublishPost(id uint) (*models.Post, error)
    GenerateSlug(title string) string
}
```

Features:
- Validation logic (title not empty, content length, etc.)
- Auto-generate slug from title if not provided
- Set published_at when status changes to "published"
- Increment view count when post is retrieved
- Pagination logic
- Business rules enforcement

#### Comment Service
Create `internal/service/comment_service.go`:

Business logic methods:
```go
type CommentService interface {
    CreateComment(postID, userID uint, content string, parentID *uint) (*models.Comment, error)
    GetComment(id uint) (*models.Comment, error)
    ListCommentsByPost(postID uint) ([]models.Comment, error)
    ListCommentsByUser(userID uint, page, pageSize int) ([]models.Comment, int64, error)
    UpdateComment(id uint, content string) (*models.Comment, error)
    DeleteComment(id uint) error
    ApproveComment(id uint) (*models.Comment, error)
    RejectComment(id uint) (*models.Comment, error)
    MarkAsSpam(id uint) (*models.Comment, error)
}
```

Features:
- Validate post exists before creating comment
- Validate parent comment exists if replying
- Content validation (not empty, max length)
- Auto-moderation logic (future: spam detection)
- Business rules enforcement

### 5. HTTP Handlers

#### Post Handler
Create `internal/handler/post_handler.go`:

Endpoints:
- `CreatePost(c *gin.Context)` - POST /api/posts
- `GetPost(c *gin.Context)` - GET /api/posts/:id
- `GetPostBySlug(c *gin.Context)` - GET /api/posts/slug/:slug
- `ListPosts(c *gin.Context)` - GET /api/posts
- `UpdatePost(c *gin.Context)` - PUT /api/posts/:id
- `DeletePost(c *gin.Context)` - DELETE /api/posts/:id
- `PublishPost(c *gin.Context)` - POST /api/posts/:id/publish
- `UnpublishPost(c *gin.Context)` - POST /api/posts/:id/unpublish

Request/Response structures:
```go
type CreatePostRequest struct {
    Title   string `json:"title" binding:"required,max=255"`
    Content string `json:"content" binding:"required"`
    Excerpt string `json:"excerpt"`
    Slug    string `json:"slug"`
}

type UpdatePostRequest struct {
    Title   *string `json:"title" binding:"omitempty,max=255"`
    Content *string `json:"content"`
    Excerpt *string `json:"excerpt"`
    Status  *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

type PostResponse struct {
    ID          uint       `json:"id"`
    Title       string     `json:"title"`
    Slug        string     `json:"slug"`
    Content     string     `json:"content"`
    Excerpt     string     `json:"excerpt"`
    AuthorID    uint       `json:"author_id"`
    Status      string     `json:"status"`
    PublishedAt *time.Time `json:"published_at"`
    ViewCount   int        `json:"view_count"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type PostListResponse struct {
    Posts      []PostResponse `json:"posts"`
    Total      int64          `json:"total"`
    Page       int            `json:"page"`
    PageSize   int            `json:"page_size"`
    TotalPages int            `json:"total_pages"`
}
```

Features:
- Proper HTTP status codes
- Error handling and validation
- Pagination support (query params: page, page_size)
- Filter support (query params: status, author_id)
- Bind JSON request bodies
- Return consistent JSON responses

#### Comment Handler
Create `internal/handler/comment_handler.go`:

Endpoints:
- `CreateComment(c *gin.Context)` - POST /api/posts/:post_id/comments
- `GetComment(c *gin.Context)` - GET /api/comments/:id
- `ListCommentsByPost(c *gin.Context)` - GET /api/posts/:post_id/comments
- `UpdateComment(c *gin.Context)` - PUT /api/comments/:id
- `DeleteComment(c *gin.Context)` - DELETE /api/comments/:id
- `ApproveComment(c *gin.Context)` - POST /api/comments/:id/approve
- `RejectComment(c *gin.Context)` - POST /api/comments/:id/reject

Request/Response structures:
```go
type CreateCommentRequest struct {
    Content  string `json:"content" binding:"required,min=1,max=1000"`
    ParentID *uint  `json:"parent_id"`
}

type UpdateCommentRequest struct {
    Content string `json:"content" binding:"required,min=1,max=1000"`
}

type CommentResponse struct {
    ID        uint       `json:"id"`
    PostID    uint       `json:"post_id"`
    UserID    uint       `json:"user_id"`
    ParentID  *uint      `json:"parent_id"`
    Content   string     `json:"content"`
    Status    string     `json:"status"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

Features:
- Proper HTTP status codes
- Nested comment support
- Status moderation
- Error handling

### 6. Route Registration

Update `cmd/server/main.go`:
```go
// Initialize repositories
postRepo := repository.NewPostRepository(database.GetDB())
commentRepo := repository.NewCommentRepository(database.GetDB())

// Initialize services
postService := service.NewPostService(postRepo)
commentService := service.NewCommentService(commentRepo, postRepo)

// Initialize handlers
postHandler := handler.NewPostHandler(postService)
commentHandler := handler.NewCommentHandler(commentService)

// Register routes
api := r.Group("/api")
{
    // Posts
    posts := api.Group("/posts")
    {
        posts.GET("", postHandler.ListPosts)
        posts.POST("", postHandler.CreatePost)
        posts.GET("/:id", postHandler.GetPost)
        posts.GET("/slug/:slug", postHandler.GetPostBySlug)
        posts.PUT("/:id", postHandler.UpdatePost)
        posts.DELETE("/:id", postHandler.DeletePost)
        posts.POST("/:id/publish", postHandler.PublishPost)
        posts.POST("/:id/unpublish", postHandler.UnpublishPost)

        // Comments for a specific post
        posts.GET("/:post_id/comments", commentHandler.ListCommentsByPost)
        posts.POST("/:post_id/comments", commentHandler.CreateComment)
    }

    // Comments (general)
    comments := api.Group("/comments")
    {
        comments.GET("/:id", commentHandler.GetComment)
        comments.PUT("/:id", commentHandler.UpdateComment)
        comments.DELETE("/:id", commentHandler.DeleteComment)
        comments.POST("/:id/approve", commentHandler.ApproveComment)
        comments.POST("/:id/reject", commentHandler.RejectComment)
    }
}
```

### 7. Validation & Error Handling

Create `internal/util/validator.go`:
- Validate post title (not empty, length)
- Validate post content (not empty)
- Validate slug format (URL-friendly)
- Validate comment content (not empty, max length)
- Validate enums (status values)

Create `internal/util/response.go`:
- Standard error response structure
- Success response helpers
- Pagination response helpers

### 8. Helper Utilities

Create `internal/util/slug.go`:
- Generate slug from title
- Ensure uniqueness
- Convert to lowercase, replace spaces with hyphens
- Remove special characters

## TESTING CHECKLIST

### Posts
- [ ] Create a new post (draft)
- [ ] Create a post with custom slug
- [ ] Auto-generate slug from title
- [ ] Retrieve post by ID
- [ ] Retrieve post by slug
- [ ] List all posts with pagination
- [ ] Filter posts by status (draft, published, archived)
- [ ] Filter posts by author
- [ ] Update post title and content
- [ ] Publish a draft post (status changes, published_at set)
- [ ] Unpublish a post
- [ ] Delete post (soft delete)
- [ ] View count increments on retrieval
- [ ] Slug uniqueness validation
- [ ] Required field validation

### Comments
- [ ] Create comment on a post
- [ ] Create nested comment (reply to comment)
- [ ] Retrieve comment by ID
- [ ] List all comments for a post
- [ ] List comments with nested structure
- [ ] Update comment content
- [ ] Delete comment (soft delete)
- [ ] Delete comment with replies (cascade)
- [ ] Approve pending comment
- [ ] Reject comment
- [ ] Mark comment as spam
- [ ] Content validation (not empty, max length)
- [ ] Verify comment belongs to existing post
- [ ] Verify parent comment exists for replies

### Edge Cases
- [ ] Create post with duplicate slug (should fail)
- [ ] Create comment on non-existent post (should fail)
- [ ] Reply to non-existent comment (should fail)
- [ ] Update non-existent post/comment (should return 404)
- [ ] Delete non-existent post/comment (should return 404)
- [ ] Empty/invalid request bodies (should return 400)
- [ ] Pagination with invalid params (should use defaults)

## DELIVERABLES

1. Database migrations (up and down)
   - `000003_create_posts.up.sql` and `.down.sql`
   - `000004_create_comments.up.sql` and `.down.sql`

2. ORM Models
   - `internal/models/post.go`
   - `internal/models/comment.go`

3. Repository Layer
   - `internal/repository/post_repository.go`
   - `internal/repository/comment_repository.go`

4. Service Layer
   - `internal/service/post_service.go`
   - `internal/service/comment_service.go`

5. HTTP Handlers
   - `internal/handler/post_handler.go`
   - `internal/handler/comment_handler.go`

6. Utilities
   - `internal/util/slug.go`
   - `internal/util/validator.go`
   - `internal/util/response.go`

7. Route registration in `cmd/server/main.go`

8. Updated documentation with API endpoints

## NOTES

### Design Decisions
- Use soft deletes for both posts and comments (via BaseModel DeletedAt)
- Cascade delete comments when post is deleted (database constraint)
- Support nested comments (one level or unlimited via parent_id)
- Post status workflow: draft → published → archived
- Comment moderation: pending → approved/rejected/spam
- Slug must be unique across all posts
- Use repository pattern for clean separation of concerns
- Use service layer for business logic

### Future Enhancements (Not in this iteration)
- Post tags/categories
- Full-text search
- Post revisions/history
- Comment reactions/likes
- Rich text editor support
- Image uploads for post content
- SEO metadata fields
- Related posts
- Comment threading depth limits
- Comment voting system

### Performance Considerations
- Add database indexes on frequently queried fields
- Use pagination for listing endpoints
- Consider caching for published posts
- Eager loading for relationships (GORM Preload)
- Add view count with debouncing (optional)

### Security Notes
- Validate all user inputs
- Sanitize HTML content (future)
- Rate limit comment creation (future)
- Verify user authorization (future with auth service)
- Prevent SQL injection (GORM handles this)
- XSS prevention for comment content (future)

### API Response Codes
- 200 OK - Successful GET, PUT
- 201 Created - Successful POST
- 204 No Content - Successful DELETE
- 400 Bad Request - Validation errors
- 404 Not Found - Resource not found
- 500 Internal Server Error - Server errors

### Migration Notes
- Run migrations in order: 000003 (posts) then 000004 (comments)
- Comments depend on posts (foreign key)
- Use CASCADE for foreign key constraints
- Test rollback migrations work correctly
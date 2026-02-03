# UI Initialization - Next.js Frontend

## Goal
Set up a Next.js frontend application with TypeScript, TailwindCSS, and API integration for the Inkstack blog platform.

## Tech Stack
- Next.js 14+ (App Router)
- TypeScript
- TailwindCSS
- Axios for API requests

## API Endpoints

### Auth Service (http://localhost:8082)
| Method | Endpoint | Payload | Description |
|--------|----------|---------|-------------|
| POST | /api/auth/register | {email, username, password} | Register new user |
| POST | /api/auth/login | {email_or_username, password} | User login |
| POST | /api/auth/logout | {refresh_token} | User logout |
| GET | /api/auth/me | - | Get current user profile |
| POST | /api/auth/change-password | {old_password, new_password} | Change password |

### API Service (http://localhost:8081)
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | /api/posts | No | List all posts |
| GET | /api/posts/:id | No | Get post details |
| POST | /api/posts | Yes | Create a new post |
| PUT | /api/posts/:id | Yes | Update a post |
| DELETE | /api/posts/:id | Yes | Delete a post |
| GET | /api/posts/:id/comments | No | List comments for a post |
| POST | /api/posts/:id/comments | Yes | Add a comment |

## Color Theme
| Role | Color | Hex |
|------|-------|-----|
| Primary background | White | #ffffff |
| Secondary background | Light gray | #f5f5f5 |
| Primary text | Dark gray | #333333 |
| Secondary text | Medium gray | #666666 |
| Accent | Ink blue | #1a365d |
| Border | Light gray | #e5e5e5 |
| Footer background | Darker gray | #d0d0d0 |

## Pages

### Login (/login)
- No header or footer
- All elements centered both vertically and horizontally
- Layout (top to bottom):
  - Title "Inkstack" — large font, ink blue, links to home
  - Spacer
  - Error notice — hidden by default; displays a red banner with an error message on login failure (e.g., "Invalid username or password")
  - Username/email input — placeholder: "Username or Email"
  - Password input — placeholder: "Password"; includes a "Forgot password?" link aligned to the upper right
  - Spacer
  - Login button — full form width, ink blue background, white text
  - Spacer
  - Footer text: "Don't have an account?" with a clickable "Register" link

### Register (/register)
- No header or footer
- All elements centered both vertically and horizontally
- Layout (top to bottom):
  - Title "Inkstack" — large font, ink blue, links to home
  - Spacer
  - Error notice — hidden by default; displays a red banner with an error message on registration failure (e.g., "Email already exists", "Password too weak")
  - Email input — placeholder: "Email"
  - Username input — placeholder: "Username"
  - Password input — placeholder: "Password (min 8 chars, 1 uppercase, 1 lowercase, 1 number, 1 special)"
  - Confirm password input — placeholder: "Confirm Password"
  - Spacer
  - Register button — full form width, ink blue background, white text
  - Spacer
  - Footer text: "Already have an account?" with a clickable "Login" link

### Visitor View (/)
- Default landing page for unauthenticated users
- Functions as a read-only subset of the Home page
- Capabilities:
  - Browse the post list
  - View post details
  - Navigate to login or register
- Layout:
  - **Header:**
    - Fixed at top with white background and a subtle bottom border
    - "Inkstack" logo on the left — links to home
    - Navigator centered, left edge aligned with post cards
    - "Login" and "Register" links on the right
  - **Main:**
    - Centered content area with a max-width of 800px
    - Post cards stacked vertically, each containing:
      - Title (clickable — navigates to post detail)
      - Author username and publish date
      - Excerpt (approximately 150 characters)
    - Cards styled with white background, subtle shadow, and rounded corners
    - Sorted by `updated_at` in descending order
    - Paginated with 10 cards per page
    - Paginator below the last card: "< Previous | Page X of Y | Next >"
  - **Footer:**
    - Horizontal divider line
    - Darker gray background (#d0d0d0)
    - Vertical padding applied
    - Content placeholder — to be added later

### Home View (/)
- Landing page for authenticated users (same route, conditional rendering)
- Layout mirrors Visitor View with these differences:
  - **Header (right side):** Displays username and a "Logout" button instead of login/register links
  - **Main:** Includes a "New Post" button above the post cards
    - Left edge aligned with post cards
    - Ink blue background, white text, rounded corners
    - Clicking navigates to the post editor

### Post Detail (/posts/[slug])
- Header and Footer consistent with Home/Visitor views
- **Main:**
  - Article card — white background, subtle shadow, max-width of 900px
  - Card contents:
    - Title — large font
    - Meta information: author username, publish date, view count
    - Horizontal divider
    - Full post content with Markdown rendering support
  - Below the article card:
    - Spacer
    - Section header: "Comments (N)"
    - If authenticated: comment textarea with a "Submit" button
    - If unauthenticated: "Login to comment" link
    - Comment list:
      - Each comment displays: username, date, content
      - Sorted by `created_at` in ascending order
      - Nested replies indented (when `parent_id` exists)

### New Post (/posts/new)
- Header and Footer consistent with Home view
- Requires authentication — redirects to login if unauthenticated
- **Main:**
  - Centered form card with a max-width of 900px
  - Form fields:
    - Title input — placeholder: "Post Title"
    - Slug input — auto-generated from title, manually editable
    - Excerpt textarea — placeholder: "Brief summary (optional)"
    - Content textarea — large area, placeholder: "Write your post content here..."
    - Status dropdown — options: "Draft" or "Published"
  - Action buttons:
    - "Cancel" — returns to home
    - "Save" — submits the form

### Edit Post (/posts/[slug]/edit)
- Layout identical to New Post
- Form pre-populated with existing post data
- Accessible only to the post author

### Profile (/profile)
- Header and Footer consistent with Home view
- Requires authentication — redirects to login if unauthenticated
- **Main:**
  - Centered card with a max-width of 600px
  - **Section: Profile Information**
    - Read-only fields: email, username
    - Editable fields: display name input, bio textarea
    - "Save Changes" button
  - Horizontal divider
  - **Section: Change Password**
    - Current password input
    - New password input
    - Confirm new password input
    - "Change Password" button
  - Horizontal divider
  - **Section: Danger Zone**
    - "Delete Account" button — red styling, triggers a confirmation dialog

### Navigator
- Horizontal navigation bar embedded within the header
- Positioned in the center of the header
- Navigation items:
  - "Posts" — navigates to home (/)
  - "Profile" — navigates to profile (/profile); visible only when authenticated

## Components

### Header
- Fixed position at the top of the viewport
- White background with a subtle bottom border (1px #e5e5e5)
- Height: 60px
- Contains: logo, navigator, authentication actions

### Footer
- Sticky positioning at the bottom
- Horizontal line at the top edge
- Darker gray background
- Minimum height: 80px
- Content placeholder for future updates

### PostCard
- Used in the post list
- White background, 8px rounded corners, subtle shadow
- Padding: 20px
- Bottom margin: 16px
- Contains: title, meta information (author, date), excerpt

### Paginator
- Centered below the post list
- Displays: previous button, "Page X of Y", next button
- Buttons disabled at boundary pages (first/last)

### CommentItem
- Used in the comment list
- Displays: username, date, content
- Nested comments indented with a left border indicator

## State Management
- Use React Context or Zustand for authentication state
- Store `access_token` and `refresh_token` in localStorage
- On application load: validate token and fetch user profile

## Authentication Flow
1. User logs in → tokens stored in localStorage
2. Each API request → `access_token` attached via Authorization header
3. On 401 response → attempt token refresh
4. If refresh fails → clear tokens and redirect to login
5. On logout → call logout API and clear localStorage

## File Structure
```
ui/
├── src/
│   ├── app/
│   │   ├── layout.tsx          # Root layout with header/footer
│   │   ├── page.tsx            # Home/Visitor page
│   │   ├── login/page.tsx
│   │   ├── register/page.tsx
│   │   ├── profile/page.tsx
│   │   └── posts/
│   │       ├── new/page.tsx
│   │       └── [slug]/
│   │           ├── page.tsx    # Post detail
│   │           └── edit/page.tsx
│   ├── components/
│   │   ├── Header.tsx
│   │   ├── Footer.tsx
│   │   ├── Navigator.tsx
│   │   ├── PostCard.tsx
│   │   ├── Paginator.tsx
│   │   ├── CommentItem.tsx
│   │   └── AuthGuard.tsx       # Route protection wrapper
│   ├── contexts/
│   │   └── AuthContext.tsx
│   ├── services/
│   │   ├── api.ts              # Axios instance configuration
│   │   ├── auth.ts             # Authentication API calls
│   │   └── posts.ts            # Posts API calls
│   ├── types/
│   │   └── index.ts            # TypeScript interfaces
│   └── lib/
│       └── utils.ts            # Helper functions
├── tailwind.config.ts
├── next.config.ts
└── package.json
```

# UI Initialization - Next.js Frontend

## GOAL
Initialize Next.js frontend application with TypeScript, TailwindCSS, and API integration for the Inkstack blog platform.

## PROJECT STRUCTURE
```
ui/
├── src/
│   ├── app/              # App router (Next.js 14+)
│   │   ├── layout.tsx    # Root layout
│   │   ├── page.tsx      # Home page (post list)
│   │   ├── posts/
│   │   │   └── [slug]/page.tsx  # Post detail page
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── components/
│   │   ├── PostCard.tsx
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   ├── services/
│   │   ├── api.ts        # Axios instance
│   │   ├── auth.ts       # Auth API calls
│   │   └── posts.ts      # Posts API calls
│   ├── types/
│   │   └── index.ts      # TypeScript types
│   └── lib/
│       └── utils.ts      # Utility functions
├── .env.local
├── .env.example
├── Dockerfile
├── next.config.js
├── tailwind.config.js
└── tsconfig.json
```

## REQUIREMENTS

### 1. Initialize Next.js Project
```bash
npx create-next-app@latest ui --typescript --tailwind --app --no-src-dir
cd ui && mkdir -p src/{components,services,types,lib}
```

### 2. Dependencies
- `axios` - HTTP client
- `react-hook-form` - Form handling
- `zustand` or `jotai` - State management (lightweight)
- `date-fns` - Date formatting

### 3. Environment Configuration
`.env.example`:
```
NEXT_PUBLIC_API_URL=http://localhost:8081/api
NEXT_PUBLIC_AUTH_URL=http://localhost:8082/api
```

### 4. API Service Layer
- `src/services/api.ts` - Axios instance with interceptors for JWT tokens
- `src/services/auth.ts` - login, register, logout, getProfile
- `src/services/posts.ts` - getPosts, getPostBySlug, createPost

### 5. Core Pages
- **Home (`/`)** - List published posts with pagination
- **Post Detail (`/posts/[slug]`)** - Display full post content
- **Login (`/login`)** - Login form with JWT token storage
- **Register (`/register`)** - Registration form

### 6. Components
- `PostCard` - Post preview card for list view
- `Header` - Navigation with login/logout
- `Footer` - Simple footer

### 7. TypeScript Types
Define interfaces matching backend models:
- `User`, `Post`, `Comment`, `LoginResponse`

### 8. Authentication
- Store JWT token in localStorage/cookies
- Add Authorization header to API requests
- Redirect to login for protected routes

### 9. Dockerfile
Multi-stage build for production optimization

## TESTING CHECKLIST
- [ ] Home page displays list of posts
- [ ] Click post navigates to detail page
- [ ] Register new user works
- [ ] Login stores JWT token
- [ ] Protected actions require authentication
- [ ] Logout clears token
- [ ] API errors are handled gracefully
- [ ] Responsive design works on mobile

## DELIVERABLES
1. Next.js project with TypeScript + TailwindCSS
2. API integration with auth and posts services
3. Authentication flow (login/register/logout)
4. Home and post detail pages
5. Dockerfile for containerization
6. Environment configuration files
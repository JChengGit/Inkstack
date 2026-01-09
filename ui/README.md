# Inkstack UI

Next.js frontend for the Inkstack blog platform.

## Tech Stack

- **Next.js 15** - React framework with App Router
- **TypeScript** - Type safety
- **TailwindCSS** - Utility-first CSS
- **Axios** - HTTP client
- **React Hook Form** - Form handling
- **Zustand** - State management
- **date-fns** - Date formatting

## Getting Started

### Prerequisites

- Node.js 20+
- Backend services running (API service on :8081, Auth service on :8082)

### Installation

```bash
# Install dependencies
npm install

# Copy environment file
cp .env.example .env.local

# Update .env.local with your API URLs
```

### Development

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build

```bash
npm run build
npm start
```

## Environment Variables

Create a `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost:8081/api
NEXT_PUBLIC_AUTH_URL=http://localhost:8082/api
```

## Project Structure

```
src/
├── app/              # App router pages
│   ├── page.tsx      # Home page (posts list)
│   ├── login/        # Login page
│   ├── register/     # Register page
│   └── posts/[slug]/ # Post detail page
├── components/       # Reusable components
│   ├── Header.tsx
│   ├── Footer.tsx
│   └── PostCard.tsx
├── services/         # API services
│   ├── api.ts        # Axios instance
│   ├── auth.ts       # Auth API
│   └── posts.ts      # Posts API
├── types/            # TypeScript types
│   └── index.ts
└── lib/              # Utilities
    ├── store.ts      # Zustand store
    └── utils.ts      # Helper functions
```

## Features

- ✅ User authentication (register/login/logout)
- ✅ Posts listing with pagination
- ✅ Post detail view
- ✅ JWT token management
- ✅ Responsive design
- ✅ Error handling

## Docker

Build and run with Docker:

```bash
# Build image
docker build -t inkstack-ui .

# Run container
docker run -p 3000:3000 \
  -e NEXT_PUBLIC_API_URL=http://api:8081/api \
  -e NEXT_PUBLIC_AUTH_URL=http://auth:8082/api \
  inkstack-ui
```

## API Integration

The UI connects to two backend services:

- **API Service** (port 8081) - Posts, comments
- **Auth Service** (port 8082) - User authentication

Authentication flow:
1. User logs in → receives JWT tokens
2. Tokens stored in localStorage
3. Axios interceptor adds token to requests
4. Protected routes require authentication

## License

MIT

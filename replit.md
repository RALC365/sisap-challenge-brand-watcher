# Brand Protection Monitor POC

## Overview
A Certificate Transparency (CT) log monitoring system for brand protection. This application monitors CT logs for certificates that match configured keywords, helping detect potential phishing or brand impersonation attempts.

## Project Structure

### Backend (Go)
- **Architecture**: Screaming Architecture / Feature Colocation
- **Location**: `backend/`
- **Entry Point**: `backend/cmd/server/main.go`

#### Core Modules:
- `internal/config/` - Application configuration
- `internal/db/` - Database connection pool and migrations
- `internal/observability/` - Logging and middleware
- `internal/scheduler/` - Job scheduler for CT log polling
- `internal/ct/` - Certificate Transparency log client
- `internal/parser/` - X.509 certificate parsing and normalization
- `internal/matcher/` - Keyword matching engine

#### Features:
- `internal/features/keywords/` - Keyword CRUD operations
- `internal/features/monitor/` - Monitor status and metrics
- `internal/features/matches/` - Certificate match results
- `internal/features/export/` - CSV export functionality

### Frontend (React + TypeScript)
- **Architecture**: Feature-driven design with strict colocation
- **Location**: `frontend/`
- **Entry Point**: `frontend/src/main.tsx`

#### Structure:
- `src/app/` - Main application component
- `src/routes/` - TanStack Router routes
- `src/providers/` - React context providers
- `src/lib/` - Utility functions and API client
- `src/components/` - Shared UI components
- `src/features/` - Feature modules (app-shell, monitor, keywords, matches, export)
- `src/styles/` - Tailwind CSS and design tokens

## Development

### Backend
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Tech Stack
- **Backend**: Go 1.25, Chi router, PostgreSQL
- **Frontend**: React 18, TypeScript, TanStack Query, Tailwind CSS, Zod validation
- **Database**: PostgreSQL 14+

## CT Log Configuration
- **URL**: https://oak.ct.letsencrypt.org/2026h2 (Oak 2026h2 - Let's Encrypt)
- **Batch Size**: 100 entries per polling cycle
- **Poll Interval**: 60 seconds

## Recent Updates
- Fixed CT log URL documentation (Oak 2026h2)
- Fixed UI layout - date range inputs no longer overlap with search
- Fixed frontend/backend parameter alignment (keyword_ids, search, start_date, end_date, limit)
- Added Export Modal (SCR-03) with filter chips, 429/500 error handling
- Added keyword highlighting in matched domain names
- Added truncated fingerprint column (hover for full hash)
- Created README.md with setup, features, design decisions, limitations

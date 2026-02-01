# Brand Protection Monitor

A Certificate Transparency (CT) log monitoring system for brand protection. This application monitors CT logs for certificates that match configured keywords, helping detect potential phishing or brand impersonation attempts.

## Setup

### Prerequisites
- Go 1.25+
- Node.js 20+
- PostgreSQL 14+

### Environment Variables
```bash
DATABASE_URL=postgresql://user:password@host:port/database
SESSION_SECRET=your-session-secret
```

### Backend Setup
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

The application will be available at:
- Frontend: http://localhost:5000
- Backend API: http://localhost:8080

## Features

### Core Functionality
- **CT Log Monitoring**: Polls Certificate Transparency logs (Argon 2025h1) for new certificates
- **Keyword Matching**: Case-insensitive matching against certificate CN and SAN fields
- **Match Deduplication**: Certificates are deduplicated by SHA256 fingerprint + keyword + matched field
- **First Seen Immutability**: `first_seen_at` timestamp never changes; `last_seen_at` updates on re-detection

### Dashboard (SCR-01)
- Real-time monitor status with idle/running/error states
- Metric cards: processed count, match count, parse errors, cycle duration
- Filter bar with keyword multi-select, date range, search, and "new only" toggle
- Sortable, paginated matches table with sticky header
- Keyword highlighting in matched domain names
- Truncated fingerprint display (hover for full hash)

### Keyword Management (SCR-02)
- Add keywords with inline validation (2-64 characters)
- Duplicate detection with 409 error handling
- Soft delete with toast notifications
- Real-time list updates via query invalidation

### Export Modal (SCR-03)
- Filter summary chips showing active filters
- Server-streamed CSV download (no client-side generation)
- Rate limiting (429) with retry guidance
- Error handling (500) with retry option
- Export audit logging to database

### API Endpoints
- `GET /health` - Health check
- `GET /monitor/status` - Monitor state and metrics
- `GET /keywords` - List keywords
- `POST /keywords` - Create keyword
- `DELETE /keywords/:id` - Soft delete keyword
- `GET /matches` - List matches with filtering/pagination
- `GET /export.csv` - Stream CSV export

## Design Decisions

### Backend Architecture
- **Screaming Architecture**: Features organized by domain (keywords, matches, monitor, export)
- **Repository Pattern**: Database access abstracted behind repositories
- **Upsert Strategy**: ON CONFLICT for deduplication, preserving first_seen_at
- **Batch Processing**: Configurable batch size (default 100 entries per cycle)
- **Rate Limiting**: Per-IP rate limiting for CSV exports

### Frontend Architecture
- **Feature Colocation**: Components, hooks, and API calls grouped by feature
- **Zod Validation**: All API responses parsed through Zod schemas before rendering
- **TanStack Query**: Server state management with automatic refetching
- **Design Tokens**: Consistent color, typography, and spacing system
- **No Persistence**: Filters reset on page reload (intentional)

### CT Log Integration
- Uses Argon 2025h1 CT log (ct.googleapis.com/logs/us1/argon2025h1)
- Fixed recent batch strategy: fetches last N entries per cycle
- Graceful error handling for CT log connection issues

## Limitations / Known Issues

1. **Single CT Log**: Currently monitors only one CT log (Argon 2025h1). Production systems should monitor multiple logs.

2. **No Authentication**: The POC does not include user authentication. Add authentication before production deployment.

3. **Limited Matching**: Only matches against CN and SAN fields. Could be extended to match issuer names, organization names, etc.

4. **No Webhook/Alerts**: Matches are only viewable in the UI. Production systems should include email/webhook notifications.

5. **Memory Usage**: Large exports may consume significant memory. Consider implementing cursor-based streaming for very large datasets.

6. **Single Instance**: No horizontal scaling support. Would need distributed locking for multi-instance deployment.

## Tech Stack

### Backend
- Go 1.25
- Chi router
- PostgreSQL with pgx driver
- Zap structured logging

### Frontend
- React 18
- TypeScript
- TanStack Query
- TanStack Router (simple path-based routing)
- Tailwind CSS
- Zod validation
- Vite build tool

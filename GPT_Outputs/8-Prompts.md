# Brand Protection Monitor (PoC)
## Replit Prompt-Driven Implementation — Source of Truth (v2)
Document date (source artifacts): 2026-02-01

> **Purpose:** A single, definitive, copy/paste prompt sequence to build the full-stack PoC in Replit from scratch.
> **Scope:** One CT Log (`https://oak.ct.letsencrypt.org/2026h2`), periodic fixed-batch ingestion, CN/SAN matching vs keywords, persist matches only, triage UI, CSV export.
> **Runtime auth:** Explicitly **excluded** (no login/JWT flows). DB may contain users table for forward compatibility.

---

# 0) Non-Negotiable Rules (Global)

## 0.1 Product and Runtime Constraints
- Single-tenant PoC. No multi-tenant abstractions.
- Monitor ONE CT Log only.
- Store ONLY matches (do not store non-matching certs).
- Deduplicate matches across cycles with unique key:
  `(cert_fingerprint_sha256, matched_keyword, matched_domain)`.
- Temporal rules:
  - `first_seen_at` is immutable.
  - `last_seen_at` updates on every reappearance (including within same run).
- Keyword deletion is **soft delete**. Deleted keywords are excluded from filters/exports; historical matches remain.
- No long global transactions per cycle. Persist via UPSERT per match (or micro-batches).

## 0.2 Engineering Constraints
- Backend: Go + Gin + pgx/v5 + Zap logging + validator.
- Frontend: React + TypeScript + Vite + TanStack Router + TanStack Query + Zustand + Axios + Zod + Tailwind.
- All backend SQL must be parameterized.
- All frontend API responses must be runtime-validated with Zod. Invalid payload => controlled error banner.
- UI is stateless across reloads (no localStorage/sessionStorage).
- Layout is flow-based (no absolute positioning). Text wraps; truncation forbidden except certificate fingerprint (if specified).

## 0.3 Deliverables (Tech Challenge)
- Working full-stack app: monitor + process + display + export CSV.
- README.md: setup, implemented features, design decisions/ambiguities, limitations/known bugs.

---

# 1) PROMPT — Create Replit Project and Folder Scaffolding

Copy/paste in Replit AI:

**Create a new workspace** named `brand-protection-monitor-poc` with two top-level folders: `backend/` and `frontend/`.

Create this backend structure (Screaming Architecture / feature colocation):

```text
backend/
  cmd/server/main.go
  internal/
    config/
      config.go
    db/
      pool.go
      migrate.go
    observability/
      logger.go
      middleware.go
    scheduler/
      scheduler.go
    ct/
      client.go
      types.go
    parser/
      x509_parser.go
      fingerprint.go
      normalize.go
    matcher/
      matcher.go
      types.go
    features/
      keywords/
        handler.go
        service.go
        repository.go
        types.go
        schemas.go
      monitor/
        handler.go
        service.go
        repository.go
        types.go
      matches/
        handler.go
        service.go
        repository.go
        types.go
      export/
        handler.go
        service.go
        repository.go
        types.go
```

Create this frontend structure (Feature-driven design with strict colocation):

```text
frontend/
  src/
    app/
      app.tsx
    main.tsx
    routes/
      __root.tsx
      index.ts
      redirect.route.tsx
      dashboard.route.tsx
      keywords.route.tsx
    providers/
      query-client-provider.tsx
      app-providers.tsx
    lib/
      axios.ts
      env.ts
      query-client.ts
      url.ts
    components/
      ui/
        button.tsx
        input.tsx
        select.tsx
        badge.tsx
        modal.tsx
        table.tsx
        toggle.tsx
        date-range.tsx
      feedback/
        error-banner.tsx
        toast.tsx
        skeleton.tsx
    features/
      app-shell/
        components/
          app-shell.tsx
          app-nav.tsx
        index.ts
      monitor/
        api/
          use-monitor-status.ts
        components/
          status-badge.tsx
          metric-card.tsx
          monitor-header.tsx
          metrics-row.tsx
        pages/
          dashboard-page.tsx
        types/
          monitor-status.schema.ts
          monitor.types.ts
        index.ts
      keywords/
        api/
          use-keywords.ts
          use-create-keyword.ts
          use-delete-keyword.ts
        components/
          keyword-form.tsx
          keyword-table.tsx
        pages/
          keywords-page.tsx
        types/
          keyword.schema.ts
          keyword.types.ts
        index.ts
      matches/
        api/
          use-matches.ts
        components/
          filter-bar.tsx
          matches-table.tsx
          highlight.tsx
        types/
          match.schema.ts
          match.types.ts
        index.ts
      export/
        api/
          use-export-csv.ts
        components/
          export-modal.tsx
          export-button.tsx
        types/
          export.types.ts
        index.ts
    styles/
      tailwind.css
      tokens.css
```

---

# 2) PROMPT — Database Setup (PostgreSQL) + Seed Singletons

Copy/paste in Replit AI:

1) Provision Postgres for the workspace.
2) Execute the FULL DDL from the PostgreSQL Source of Truth:
   - Extensions: `uuid-ossp`, `pgcrypto`, `pg_trgm`
   - Enums: `keyword_status`, `monitor_state_enum`, `monitor_run_state_enum`, `matched_field_enum`
   - Tables: `users`, `keywords`, `monitor_config`, `monitor_state`, `monitor_run`, `matched_certificates`, `exports`
   - Triggers: `trigger_set_updated_at`, `trigger_bump_version_number`, and per-table triggers
   - Indexes: including trigram indexes for substring search and unique indexes for dedupe

3) Seed singleton rows (must exist before backend starts):
- `monitor_config` row with `id=1` and:
  - `ct_log_base_url = 'https://oak.ct.letsencrypt.org/2026h2'`
  - `poll_interval_seconds = 60`
  - `batch_size = 100`
  - `ct_connect_timeout_ms = 2000`
  - `ct_read_timeout_ms = 5000`
- `monitor_state` row with `id=1`, `state='idle'`

---

# 3) PROMPT — Backend Bootstrap (Go + Gin + pgx)

Copy/paste in Replit AI:

## 3.1 Initialize Go module and deps
In `backend/`:

```bash
go mod init brand-protection-monitor
go get github.com/gin-gonic/gin
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool
go get go.uber.org/zap
go get github.com/go-playground/validator/v10
```

## 3.2 Implement config loader (`internal/config/config.go`)
- Read env:
  - `DATABASE_URL`
  - `PORT` (default 8080)
- Validate required envs at startup; fail-fast with clear logs.

## 3.3 Implement DB pool singleton (`internal/db/pool.go`)
- Create pgxpool with sane defaults.
- Expose `GetPool() *pgxpool.Pool` and `ClosePool()`.

## 3.4 Implement logger + middleware (`internal/observability/*`)
- Zap JSON logger singleton.
- Gin middleware:
  - request_id correlation (generate UUID if missing)
  - structured access log per request
  - panic recovery
  - basic rate limiting middleware (per-IP) for:
    - /matches (search heavy)
    - /export.csv (stricter)
- All logs include `request_id` and when relevant `run_id`.

## 3.5 Implement health endpoints
- `GET /healthz`: returns 200 if process alive.
- `GET /readyz`: checks DB connectivity and ability to read `monitor_state` id=1. If fails => 500 DB_UNAVAILABLE.

## 3.6 main.go composition (`cmd/server/main.go`)
- Load config
- Init logger
- Init DB pool
- Init router and register all feature routes
- Start scheduler (monitor loop) in background goroutine
- Start HTTP server

---

# 4) PROMPT — REST API Contracts (Backend MUST MATCH)

Implement the exact API surface (JSON unless CSV):

## 4.1 GET /monitor/status
Response fields:
- `state`: "idle" | "running" | "error"
- `last_run_at`: timestamp|null
- `last_success_at`: timestamp|null
- `last_error_code`: string|null
- `last_error_message`: string|null
- `metrics_last_run`: object|null with:
  - `processed_count`
  - `match_count`
  - `parse_error_count`
  - `duration_ms`
  - `ct_latency_ms`
  - `db_latency_ms`
Errors:
- 500 DB_UNAVAILABLE

## 4.2 Keywords
### GET /keywords
Query:
- `q?` substring search
- `page?` int
- `page_size?` one of 10|25|50
Response:
- `items`: array
- `total`: int
Errors:
- 400 INVALID_QUERY
- 500 DB_ERROR

### POST /keywords
Body: `{ "value": string }`
Server rules:
- Reject missing key or non-string types
- Trim value
- Reject empty/whitespace-only
- Enforce length 1..64
- normalized_value = lower(trim(value))
- Enforce case-insensitive dedupe via unique index on normalized_value where is_deleted=false
Response: created keyword (keyword_id, value, normalized_value, status, created_at)
Errors:
- 400 VALIDATION_ERROR
- 409 DUPLICATE_KEYWORD
- 500 DB_ERROR

### DELETE /keywords/{keyword_id}
- Path param must be numeric
- Soft delete: set status='inactive', is_deleted=true, deleted_at=NOW()
Response: `{ "ok": true }`
Errors:
- 400 INVALID_PATH_PARAM
- 404 NOT_FOUND
- 500 DB_ERROR

## 4.3 GET /matches
Query params:
- `page`, `page_size` (10|25|50)
- `keyword` (recommended: normalized keyword value)
- `q` (substring across domain OR issuer)
- `issuer` (substring issuer only)
- `date_from`, `date_to` (filter by first_seen_at inclusive)
- `new_only` boolean
- `sort`: `first_seen_desc` (default) | `last_seen_desc` | `domain_asc`
Response:
- `{ items: [...], total: number }`
Errors:
- 400 INVALID_QUERY
- 500 DB_ERROR
Rate limit applies.

## 4.4 GET /export.csv
- Same filters as /matches
- Content-Type: text/csv
- Must be streaming (no buffering entire dataset)
- Must record `exports` row with filters snapshot JSONB
Errors:
- 429 RATE_LIMITED
- 500 EXPORT_ERROR (must occur before headers if possible)

---

# 5) PROMPT — Module: Keyword Management (Backend Feature Implementation)

Implement `internal/features/keywords` with:

## Files and Responsibilities
- `handler.go`: Gin handlers (request parsing, status codes, response envelope)
- `service.go`: validations and domain logic (trim, normalize, length checks, error mapping)
- `repository.go`: pgx queries (parameterized), transactions short-lived
- `schemas.go`: request/response structs + validator tags (Go-side)
- `types.go`: shared error codes, DTOs

## Required Flow (POST /keywords)
1. Parse JSON body.
2. Validate type and presence of `value`.
3. Trim and validate non-empty.
4. Enforce 1..64 length.
5. Compute normalized_value.
6. Attempt INSERT.
7. If unique violation => 409 DUPLICATE_KEYWORD.
8. Return created payload.

## Required Flow (DELETE /keywords/{keyword_id})
1. Parse path param and validate numeric.
2. Update row where is_deleted=false.
3. If 0 rows => 404 NOT_FOUND.

## Required Flow (GET /keywords)
- Apply `is_deleted=false`
- Optional substring search on normalized_value
- Order by keyword_id desc
- Return paginated items and total count

---

# 6) PROMPT — Module: CT Client + Parsing + Matching + Persistence

## 6.1 CT Client (`internal/ct`)
Implement:
- `GET {ct_log_base_url}/ct/v1/get-sth`
- `GET {ct_log_base_url}/ct/v1/get-entries?start=X&end=Y`

Rules:
- Apply connect/read timeouts from `monitor_config`.
- Parse JSON strictly; if invalid => fatal cycle error.
- Support chunking get-entries requests if needed.

## 6.2 Range Calculation
- end = tree_size - 1
- start = max(0, end - batch_size)

Persist start/end into `monitor_run`.

## 6.3 X.509 Parsing (`internal/parser`)
For each CT entry:
- Extract leaf certificate from entry.
- Extract CN + SAN DNS names.
- Normalize domains: lowercase, strip trailing dot.
- Compute SHA-256 fingerprint hex.

Parse errors are non-fatal per certificate; increment parse_error_count and continue.

## 6.4 Matching (`internal/matcher`)
- Load active keywords once per cycle (is_deleted=false, status=active).
- Case-insensitive substring match on CN and each SAN.
- Record matched_field: cn | san | both.

## 6.5 UPSERT Pattern
Use ON CONFLICT on (cert_fingerprint_sha256, matched_keyword, matched_domain) WHERE is_deleted=false:
- Insert: first_seen_at=NOW(), last_seen_at=NOW(), last_run_id
- Update: last_seen_at=NOW(), last_run_id
Never mutate first_seen_at.

## 6.6 Monitor State Machine + Single-Flight (`internal/scheduler`)
- Lock monitor_state row FOR UPDATE
- If not idle: exit
- Set running + clear error fields
- Run cycle; persist monitor_run + metrics
- On success: set idle + last_success_at
- On fatal failure: set state=error + last_error_code/message
Partial persistence remains.

---

# 7) PROMPT — Module: Matches & Triage (Backend Feature Implementation)

Implement `internal/features/matches`:

- Validate query params (page_size, sort, dates). Invalid => 400 INVALID_QUERY.
- Build SQL filters:
  - keyword match on matched_keyword (normalized)
  - q across domain or issuer
  - issuer only filter
  - date range on first_seen_at
  - new_only derived from last successful run window (deterministic)
- Pagination: LIMIT/OFFSET + total count query.
- Sorting: first_seen_desc (default), last_seen_desc, domain_asc.

---

# 8) PROMPT — Module: CSV Export (Backend Feature Implementation)

Implement `internal/features/export`:

- Stream CSV from DB rows; no buffering.
- RFC4180 compliant escaping/quoting.
- Insert/update `exports` row with filters JSONB and rows_exported.
- Rate limit: if exceeded => 429 RATE_LIMITED.

---

# 9) PROMPT — Frontend Bootstrap (UI Tokens + Libraries)

Install:
- TanStack Router, TanStack Query, Zustand, Axios, Zod, Tailwind CSS

Implement canonical tokens exactly:
- color.primary #2563EB
- color.success #16A34A
- color.warning #F59E0B
- color.error #DC2626
- surface.page #F9FAFB
- surface.card #FFFFFF
- text.primary #111827
- text.muted #6B7280

Typography: Inter; H1/H2/H3/body sizes as specified; all text wraps.

---

# 10) PROMPT — Screen SCR-01 Dashboard (States + Layout)

Compose vertically:
1) App Bar (sticky): name + StatusBadge + last_run_at + Export CTA
2) Metric cards row: processed_count, match_count, parse_error_count, cycle_duration_ms
3) Filter bar: keyword multi-select + date range + text search + new_only toggle (300ms debounce)
4) Matches table: server-side pagination, sorting, sticky table header, long text wrap

States:
- Loading: skeleton cards + skeleton rows
- Empty no keywords: CTA to /keywords
- Empty no matches: empty educational state
- Error: persistent banner with last_error_code/message; table and export remain usable

---

# 11) PROMPT — Screen SCR-02 Keywords Management

- Add keyword (inline validation + show 409 DUPLICATE_KEYWORD inline)
- List keywords
- Soft delete keyword (toast on success, 404 toast, 500 banner)
- After delete: keyword removed from dashboard filter options

---

# 12) PROMPT — Screen SCR-03 Export Modal

- Must NOT navigate away from /dashboard.
- Show filter chips summary.
- Disable CTA while exporting.
- Handle:
  - 429 RATE_LIMITED => guidance
  - 500 EXPORT_ERROR => modal error state
- Download is server streamed only (no client-side CSV generation).

---

# 13) FINAL ACCEPTANCE CHECKLIST

Backend:
- [ ] CT Log URL configured as required
- [ ] Fixed recent batch strategy (e.g., last 100 entries) per cycle
- [ ] Keyword persistence and soft delete
- [ ] Matches persisted and deduplicated; first_seen immutable; last_seen updates
- [ ] /monitor/status metrics reflect last run
- [ ] /matches filters and sorts deterministic
- [ ] /export.csv streaming + rate limit + export audit row

Frontend:
- [ ] UI uses canonical tokens and required component library only
- [ ] Zod parses every API response before rendering
- [ ] Dashboard + Keywords + Export Modal implemented
- [ ] Highlighting is obvious and non-truncating (except fingerprint if allowed)
- [ ] No persistence across reload

README:
- [ ] setup/run instructions
- [ ] implemented features list
- [ ] design decisions
- [ ] limitations/known bugs

---
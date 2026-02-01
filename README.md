# Brand Protection Monitor

A Certificate Transparency (CT) log monitoring system for brand protection. This application monitors CT logs for certificates that match configured keywords, helping detect potential phishing or brand impersonation attempts.

**Live Demo:** [https://brand-watcher-challenge.replit.app](https://brand-watcher-challenge.replit.app)

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Installation & Local Setup](#2-installation--local-setup)
3. [Methodology: The AI Machine](#3-methodology-the-ai-machine)
4. [Technical Decisions](#4-technical-decisions)
5. [Features](#5-features)
6. [API Reference](#6-api-reference)
7. [Limitations](#7-limitations)
8. [Tech Stack](#8-tech-stack)

---

## 1. Project Overview

The **Brand Protection Monitor** is a Proof of Concept (PoC) system that monitors Certificate Transparency (CT) logs in real-time to detect potential brand impersonation or phishing attempts.

### What it does:

- **Polls CT logs** every 60 seconds for new SSL/TLS certificates
- **Matches keywords** against certificate Common Name (CN) and Subject Alternative Names (SAN)
- **Deduplicates matches** by certificate fingerprint + keyword + matched field
- **Provides a dashboard** to view, filter, and export matches
- **Exports to CSV** for further analysis or reporting

### Use Case:

Organizations can configure keywords (brand names, product names, etc.) and the system will alert them when certificates are issued containing those keywords. This helps identify potential phishing domains before they can be used in attacks.

---

## 2. Installation & Local Setup

### Prerequisites

- Go 1.25+
- Node.js 20+
- PostgreSQL 14+

### Step 1: Clone the Repository

```bash
git clone https://github.com/your-username/brand-protection-monitor.git
cd brand-protection-monitor
```

### Step 2: Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
DATABASE_URL=postgresql://user:password@localhost:5432/brand_monitor
SESSION_SECRET=your-session-secret-here
```

### Step 3: Set Up PostgreSQL Database

```bash
# Create the database
createdb brand_monitor

# The application will automatically run migrations on startup
```

### Step 4: Start the Backend

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

The backend API will be available at `http://localhost:8080`

### Step 5: Start the Frontend

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5000`

### Quick Start (Both Services)

```bash
# Terminal 1 - Backend
cd backend && go run cmd/server/main.go

# Terminal 2 - Frontend
cd frontend && npm run dev
```

---

## 3. Methodology: The AI Machine

This section describes **"The AI Machine"**, an **AI-first engineering process** used to design and implement this PoC.

### 3.1 Purpose

The goal of this methodology is to:

- Explicitly document **how AI was used**, as allowed and encouraged by the challenge
- Explain **each GPT agent involved**, its objective, inputs, and outputs
- Describe the **technical and architectural decisions** taken, including trade-offs and rationale
- Demonstrate **ownership, understanding, and human validation** of all AI-generated artifacts

### 3.2 Core Principles

**The AI Machine** is an **AI-first, human-validated delivery pipeline**:

1. **AI as the primary generator** of specifications, schemas, and architecture
2. **Human oversight at every boundary** (validation, corrections, constraints)
3. **Single-source-of-truth documents** produced per domain
4. **Zero ambiguity tolerance** — all assumptions must be resolved explicitly

### 3.3 Conceptual Flow

```
Problem Definition (Tech Challenge)
        ↓
Context Normalization (AI)
        ↓
Domain-Specific GPT Agents
        ↓
Formal Source-of-Truth Documents (PDF)
        ↓
Human Review & Constraint Enforcement
        ↓
Code Generation (Replit)
        ↓
Final PoC
```

### 3.4 GPT Agents Used

The AI Machine was executed through a **strict, sequential chain of specialized GPT agents**:

| Order | Agent | Responsibility | Output |
|-------|-------|----------------|--------|
| 1 | **Normalizer** | Parse requirements, eliminate ambiguity, define scope boundaries | Normalized problem definition |
| 2 | **PRD Generator** | Translate requirements into product language, define features and constraints | PRD (PDF) |
| 3 | **User Stories** | Convert PRD into explicit, testable user stories with acceptance criteria | User Stories (PDF) |
| 4A | **UX/UI Architect** | Define visual tokens, UI components, screen layouts, and UX behavior | UX/UI Specification (PDF) |
| 4B | **Database Architect** | Design PostgreSQL schema with deduplication and auditability | Database Schema (PDF) |
| 5A | **Tech Lead** | Define CT log polling strategy, API contracts, error handling | Technical Bible (PDF) |
| 5B | **Infra/Frontend** | Define frontend structure, routing, state management | Frontend Architecture (PDF) |
| 6 | **Prompt Generator** | Generate execution-ready prompts for Replit | Replit Prompts |

### 3.5 Human Validation

All AI outputs were:

- Reviewed line-by-line
- Constrained with explicit rules
- Adjusted when misaligned with requirements
- Fully understood and explainable by the author

**AI was used as a multiplier, not a replacement.**

---

## 4. Technical Decisions

### 4.1 Backend: Go

| Decision | Rationale |
|----------|-----------|
| Go language | Required by challenge; strong concurrency model for CT polling; deterministic performance |
| Screaming Architecture | Features organized by domain (keywords, matches, monitor, export) for clarity |
| Chi Router | Lightweight, idiomatic Go HTTP router |
| Repository Pattern | Database access abstracted for testability |
| Upsert Strategy | ON CONFLICT for deduplication, preserving first_seen_at |

### 4.2 Database: PostgreSQL

| Decision | Rationale |
|----------|-----------|
| PostgreSQL | Required by challenge; strong relational guarantees |
| GIN/Trigram Indexes | Efficient substring matching for keyword search |
| Soft Deletes | Keywords are never hard deleted; preserves match history |
| Enum Types | Type safety for matched_field (cn, san, both) |

### 4.3 Frontend: React + TypeScript

| Decision | Rationale |
|----------|-----------|
| React 18 + TypeScript | Required by challenge; type safety reduces bugs |
| TanStack Query | Server state management with automatic refetching |
| Tailwind CSS | Fast iteration with utility-first CSS |
| Zod Validation | All API responses parsed through schemas before rendering |
| Feature Colocation | Components, hooks, and API calls grouped by feature |

### 4.4 CT Log Integration

| Decision | Rationale |
|----------|-----------|
| Oak 2026h2 Log | Let's Encrypt CT log with high certificate volume |
| Fixed Recent Batch | Fetches last N entries per cycle (configurable) |
| 60-second Poll Interval | Balance between freshness and API load |
| Match-only Storage | Only certificates with keyword matches are persisted |

### 4.5 Why AI-First Development

| Benefit | Description |
|---------|-------------|
| Faster Exploration | AI explores complex design spaces quickly |
| Explicit Documentation | Forces all decisions to be documented |
| Reduced Assumptions | Hidden assumptions are surfaced and resolved |
| Human Review | Ensures correctness and alignment with requirements |

---

## 5. Features

### 5.1 Dashboard (SCR-01)

- Real-time monitor status with idle/running/error states
- Metric cards: processed count, match count, parse errors, cycle duration
- Filter bar with keyword multi-select, date range, search, and "new only" toggle
- Sortable, paginated matches table with sticky header
- Keyword highlighting in matched domain names
- Truncated fingerprint display (hover for full hash)
- Countdown timer showing time until next polling cycle

### 5.2 Keyword Management (SCR-02)

- Add keywords with inline validation (2-64 characters)
- Duplicate detection with 409 error handling
- Soft delete with toast notifications
- Real-time list updates via query invalidation

### 5.3 Export Modal (SCR-03)

- Filter summary chips showing active filters
- Server-streamed CSV download (no client-side generation)
- Rate limiting (429) with retry guidance
- Error handling (500) with retry option
- Export audit logging to database

### 5.4 CT Log Monitoring

- Connects to Oak 2026h2 CT log (Let's Encrypt)
- Polls every 60 seconds for new certificates
- Parses X.509 certificates for CN and SAN fields
- Case-insensitive keyword matching
- Deduplication by certificate SHA256 + keyword + matched field

---

## 6. API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/monitor/status` | Monitor state and metrics |
| GET | `/keywords` | List all keywords |
| POST | `/keywords` | Create a new keyword |
| DELETE | `/keywords/:id` | Soft delete a keyword |
| GET | `/matches` | List matches with filtering/pagination |
| GET | `/export.csv` | Stream CSV export |

### Query Parameters for `/matches`

| Parameter | Type | Description |
|-----------|------|-------------|
| `keyword_ids` | string[] | Filter by keyword IDs |
| `search` | string | Search in matched domain |
| `start_date` | ISO date | Filter by first_seen_at >= date |
| `end_date` | ISO date | Filter by first_seen_at <= date |
| `new_only` | boolean | Only show recent matches |
| `sort_by` | string | Sort field (first_seen_at, keyword, etc.) |
| `sort_order` | string | asc or desc |
| `page` | number | Page number (1-indexed) |
| `limit` | number | Items per page |

---

## 7. Limitations

These limitations are **intentional** and aligned with the PoC scope:

| Limitation | Description |
|------------|-------------|
| Single CT Log | Currently monitors only Oak 2026h2. Production systems should monitor multiple logs |
| No Authentication | The PoC does not include user authentication |
| No Real-time Streaming | Uses polling instead of WebSocket/SSE |
| No Webhook/Alerts | Matches are only viewable in the UI; no email/webhook notifications |
| No Takedown Workflows | No integration with domain registrars or hosting providers |
| Single Instance | No horizontal scaling support; would need distributed locking |
| Memory Usage | Large exports may consume significant memory |

---

## 8. Tech Stack

### Backend

| Technology | Purpose |
|------------|---------|
| Go 1.25 | Programming language |
| Chi | HTTP router |
| pgx | PostgreSQL driver |
| Zap | Structured logging |

### Frontend

| Technology | Purpose |
|------------|---------|
| React 18 | UI framework |
| TypeScript | Type safety |
| TanStack Query | Server state management |
| TanStack Router | Client-side routing |
| Tailwind CSS | Styling |
| Zod | Schema validation |
| Vite | Build tool |

### Infrastructure

| Technology | Purpose |
|------------|---------|
| PostgreSQL 14+ | Database |
| Replit | Hosting and deployment |

---

## Author

**Richardson Cárcamo**

Full Stack Engineering Challenge — Brand Protection Monitor (PoC)

---

*End of README*

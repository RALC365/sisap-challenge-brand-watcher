# Brand Protection Monitor

A Certificate Transparency (CT) log monitoring system for brand protection. This application monitors CT logs for certificates that match configured keywords, helping detect potential phishing or brand impersonation attempts.

**Live Demo:** [https://brand-watcher-challenge.replit.app](https://brand-watcher-challenge.replit.app)

**Repository:** [https://github.com/RALC365/sisap-challenge-brand-watcher](https://github.com/RALC365/sisap-challenge-brand-watcher)

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Installation & Local Setup](#2-installation--local-setup)
3. [Methodology: The AI Machine](#3-methodology-the-ai-machine)
   - [3.1 Purpose](#31-purpose)
   - [3.2 Alignment with Tech Challenge Requirements](#32-alignment-with-tech-challenge-requirements)
   - [3.3 Core Principles](#33-core-principles)
   - [3.4 Conceptual Flow](#34-conceptual-flow)
   - [3.5 GPT Agents Used](#35-gpt-agents-used)
   - [3.6 Human Validation and Ownership](#36-human-validation-and-ownership)
4. [Technical Decisions](#4-technical-decisions)
   - [4.1 Why Go for the Backend](#41-why-go-for-the-backend)
   - [4.2 Why PostgreSQL](#42-why-postgresql)
   - [4.3 Why React + TypeScript + Tailwind](#43-why-react--typescript--tailwind)
   - [4.4 Why Not Store All Certificates](#44-why-not-store-all-certificates)
   - [4.5 Why AI-First Instead of Manual Design](#45-why-ai-first-instead-of-manual-design)
5. [Features](#5-features)
6. [API Reference](#6-api-reference)
7. [Limitations](#7-limitations)
8. [Tech Stack](#8-tech-stack)
9. [Conclusion](#9-conclusion)

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
git clone https://github.com/RALC365/sisap-challenge-brand-watcher.git
cd sisap-challenge-brand-watcher
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
go run ./cmd/server/main.go
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
cd backend && go run ./cmd/server/main.go

# Terminal 2 - Frontend
cd frontend && npm run dev
```

---

## 3. Methodology: The AI Machine

This section describes **"The AI Machine"**, an **AI-first engineering process** used to design and implement the *Brand Protection Monitor (PoC)* required by the SISAP Tech Challenge.

### 3.1 Purpose

The goal of this documentation is to:

- Explicitly document **how AI was used**, as allowed and encouraged by the challenge
- Explain **each GPT agent involved**, its objective, inputs, and outputs
- Describe the **technical and architectural decisions** taken along the way, including trade-offs and rationale
- Demonstrate **ownership, understanding, and human validation** of all AI-generated artifacts

This document is intended to complement the source code and the generated PDFs, and to satisfy the **"Communication"** and **"Design Decisions / Ambiguities"** evaluation criteria of the challenge.

### 3.2 Alignment with Tech Challenge Requirements

The *Tech Challenge* explicitly states that:

- **AI usage is permitted and encouraged**
- The candidate must **understand, validate, and explain** any AI-generated output
- A detailed `README.md` describing **design decisions, trade-offs, and limitations** is required

> This process was intentionally designed to exceed those expectations by making AI usage **explicit, structured, auditable, and human-validated**.

### 3.3 Core Principles

**The AI Machine** is an **AI-first, human-validated delivery pipeline**:

1. **AI as the primary generator** of specifications, schemas, and architecture
2. **Human oversight at every boundary** (validation, corrections, constraints)
3. **Single-source-of-truth documents** produced per domain
4. **Zero ambiguity tolerance** — all assumptions must be resolved explicitly

### 3.4 Conceptual Flow

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

### 3.5 GPT Agents Used

The AI Machine was executed through a **strict, sequential chain of specialized GPT agents**. Each GPT consumes the validated output of the previous one. Skipping or reordering agents is explicitly forbidden in this process.

> **Note:** All generated Source-of-Truth documents (PDFs) from each GPT agent can be found in the [`/GPT_Outputs`](./GPT_Outputs) folder.

#### Execution Order

```
1. Normalizer
2. PRD
3. User Stories
4A. UX/UI
4B. Database
5A. Tech Lead
5B. Infra
6. Prompt Generator
```

Each GPT has a **non-overlapping responsibility**, clear inputs, and a concrete output artifact (PDF or prompt set).

---

#### GPT 1 — Normalizer

**Objective:** Normalize and lock the problem space before any design or implementation decisions.

**Primary Responsibilities:**
- Parse the Tech Challenge requirements
- Eliminate ambiguity
- Define scope boundaries (in-scope / out-of-scope)
- Produce resolved assumptions as binding constraints

**Inputs:** Tech Challenge PDF

**Outputs:** Normalized problem definition, Explicit assumptions, Scope lock (PoC boundaries)

**Why this GPT exists:** This agent prevents all downstream GPTs from making implicit assumptions or inventing requirements.

---

#### GPT 2 — PRD Generator

**Objective:** Generate a **formal Product Requirements Document (PRD)** from the normalized context.

**Primary Responsibilities:**
- Translate normalized requirements into product language
- Define features, non-goals, and constraints
- Establish acceptance criteria at the product level

**Inputs:** Normalized context (GPT 1 output)

**Outputs:** PRD — Source of Truth (PDF)

**Why this GPT exists:** Separates *problem understanding* from *solution design* and ensures product intent is explicit.

---

#### GPT 3 — User Stories Generator

**Objective:** Convert the PRD into **explicit, testable user stories**.

**Primary Responsibilities:**
- Produce user stories with acceptance criteria
- Define happy paths and edge cases
- Ensure full traceability back to the PRD

**Inputs:** PRD (GPT 2 output)

**Outputs:** User Stories — Source of Truth (PDF)

**Why this GPT exists:** Prevents feature gaps and provides a bridge between product intent and technical execution.

---

#### GPT 4A — UX/UI Architect

**Objective:** Define a **binding UX/UI specification** with zero visual or interaction ambiguity.

**Primary Responsibilities:**
- Define visual tokens (colors, typography)
- Define allowed UI components
- Define screen layouts and UI states
- Specify UX behavior for loading, empty, and error states

**Inputs:** PRD, User Stories

**Outputs:** UX/UI Source of Truth (PDF)

**Why this GPT exists:** Locks the frontend experience before any code is written.

---

#### GPT 4B — Database Architect

**Objective:** Design a **production-grade PostgreSQL schema** aligned with product and UX semantics.

**Primary Responsibilities:**
- Model persistence rules
- Define soft deletes and deduplication
- Ensure auditability and observability

**Inputs:** PRD, User Stories

**Outputs:** PostgreSQL Source of Truth (PDF)

**Why this GPT exists:** Ensures data integrity and long-term correctness before backend logic is implemented.

---

#### GPT 5A — Tech Lead (Backend Orchestration)

**Objective:** Define **how the system actually runs** at the backend level.

**Primary Responsibilities:**
- Define CT log polling strategy
- Define scheduling and concurrency rules
- Define REST API contracts
- Classify and handle error conditions

**Inputs:** PRD, User Stories, Database Source of Truth

**Outputs:** Technical Bible / Backend Source of Truth (PDF)

**Why this GPT exists:** Bridges architecture and implementation with explicit execution rules.

---

#### GPT 5B — Infra / Frontend Architecture

**Objective:** Define the **frontend and infrastructure architecture** needed to implement the PoC.

**Primary Responsibilities:**
- Define frontend project structure
- Define routing and state management rules
- Define API consumption patterns

**Inputs:** UX/UI Source of Truth, Technical Bible

**Outputs:** Frontend Architecture Source of Truth (PDF)

**Why this GPT exists:** Ensures the frontend is scalable, testable, and aligned with backend contracts.

---

#### GPT 6 — Prompt Generator

**Objective:** Generate **high-quality, execution-ready prompts** for Replit.

**Primary Responsibilities:**
- Translate all Source-of-Truth documents into prompts
- Ensure prompts are deterministic and modular
- Optimize prompts for AI-assisted coding

**Inputs:** All previous GPT outputs

**Outputs:** Replit-ready prompts (Markdown / text)

**Why this GPT exists:** Transforms design and architecture into executable AI instructions.

---

### 3.6 Human Validation and Ownership

All AI outputs were:

- **Reviewed line-by-line**
- **Constrained with explicit rules**
- **Adjusted when misaligned with requirements**
- **Fully understood and explainable by the author**

**AI was used as a multiplier, not a replacement.**

---

## 4. Technical Decisions

### 4.1 Why Go for the Backend

| Rationale |
|-----------|
| Required by the challenge |
| Strong concurrency model for CT polling |
| Deterministic performance |

**Additional Decisions:**

| Decision | Rationale |
|----------|-----------|
| Screaming Architecture | Features organized by domain (keywords, matches, monitor, export) for clarity |
| Chi Router | Lightweight, idiomatic Go HTTP router |
| Repository Pattern | Database access abstracted for testability |
| Upsert Strategy | ON CONFLICT for deduplication, preserving first_seen_at |

### 4.2 Why PostgreSQL

| Rationale |
|-----------|
| Required by the challenge |
| Strong relational guarantees |
| Advanced indexing (GIN, trigram) |

**Additional Decisions:**

| Decision | Rationale |
|----------|-----------|
| GIN/Trigram Indexes | Efficient substring matching for keyword search |
| Soft Deletes | Keywords are never hard deleted; preserves match history |
| Enum Types | Type safety for matched_field (cn, san, both) |

### 4.3 Why React + TypeScript + Tailwind

| Rationale |
|-----------|
| Required by the challenge |
| Fast iteration for PoC |
| Strong type safety and UI consistency |

**Additional Decisions:**

| Decision | Rationale |
|----------|-----------|
| TanStack Query | Server state management with automatic refetching |
| Zod Validation | All API responses parsed through schemas before rendering |
| Feature Colocation | Components, hooks, and API calls grouped by feature |

### 4.4 Why Not Store All Certificates

| Rationale |
|-----------|
| Explicitly allowed by PoC scope |
| Reduces storage and complexity |
| Focuses system on actionable signal (matches) |

Only certificates that match configured keywords are persisted. This is a deliberate design decision to optimize for the PoC use case.

### 4.5 Why AI-First Instead of Manual Design

| Benefit | Description |
|---------|-------------|
| Faster Exploration | AI explores complex design spaces quickly |
| Explicit Documentation | Forces all decisions to be documented |
| Reduced Assumptions | Hidden assumptions are surfaced and resolved |

Human review ensured correctness and alignment.

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
- "How It Works" instructional section

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
| GET | `/` | Root health check |
| GET | `/health` | Health check |
| GET | `/healthz` | Kubernetes-style health check |
| GET | `/readyz` | Readiness check with DB verification |
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

All limitations were **intentional** and aligned with the challenge scope.

---

## 8. Tech Stack

### Backend

| Technology | Purpose |
|------------|---------|
| Go 1.25 | Programming language |
| Gin | HTTP router |
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

## 9. Conclusion

The **AI Machine** demonstrates:

- **Strong architectural thinking**
- **Responsible AI usage**
- **Clear communication**
- **Production-grade design under time constraints**

This approach reflects how modern engineering teams can safely and effectively integrate AI into real-world delivery pipelines.

---

## Author

**Richardson Cárcamo**

Full Stack Engineering Challenge — Brand Protection Monitor (PoC)

Challenge Owner: SISAP (Recruitment Process)

---

*End of README*

# 🎵 MaqamDetector

> **Pendeteksi Maqam Musik Arab & Timur Tengah** — Identifikasi maqam secara otomatis dari audio YouTube, file upload, atau rekaman suara. Dibangun untuk komunitas banjari Indonesia.

---

## Overview

MaqamDetector membantu musisi, arranger, dan penggemar musik banjari untuk mengidentifikasi maqam (tangga nada modal) dari lagu-lagu Arab dan Timur Tengah. Output berupa nama maqam, confidence score, kandidat alternatif, dan penjelasan lengkap dalam Bahasa Indonesia.

### Maqam yang Didukung (v1)

| ID | Nama Latin | Nama Arab |
|----|-----------|-----------|
| `hijaz` | Hijaz | حجاز |
| `rast` | Rast | راست |
| `bayati` | Bayati | بياتي |
| `nahawand` | Nahawand | نهاوند |
| `kurd` | Kurd | كرد |
| `saba` | Saba | صبا |
| `ajam` | Ajam | عجم |
| `jiharkah` | Jiharkah | جهاركاه |

---

## Architecture

```
maqam-detector/
├── apps/
│   ├── platform/       ← React SPA (TanStack Router + Vite + Tailwind v4)
│   ├── api/            ← Go REST API (Echo + GORM + PostgreSQL)
│   └── analyzer/       ← Python FastAPI (CREPE + librosa + yt-dlp)
├── docs/               ← Dokumentasi produk (PRD, API Contract, dll)
├── docker-compose.yml  ← Orchestration semua service
├── nx.json             ← Nx build orchestrator
├── pnpm-workspace.yaml
└── .env
```

### Data Flow

```
Audio Input (YouTube URL / File / Mikrofon)
    ↓
Golang API → validasi, session management
    ↓
Python Analyzer → yt-dlp + CREPE pitch extraction + maqam matching
    ↓
Claude API → penjelasan kontekstual Bahasa Indonesia
    ↓
Output: Nama Maqam + Confidence Score + Penjelasan
```

---

## Tech Stack

### Backend API (`apps/api`)

| Kategori       | Teknologi |
|----------------|-----------|
| Language       | Go 1.25 |
| Framework      | [Echo v4](https://echo.labstack.com/) |
| ORM            | [GORM](https://gorm.io/) + PostgreSQL |
| Validation     | [go-playground/validator](https://github.com/go-playground/validator) |
| API Docs       | [Swagger](https://github.com/swaggo/swag) (auto-generated) |
| Hot Reload     | [Air](https://github.com/air-verse/air) |
| Logging        | [Logrus](https://github.com/sirupsen/logrus) |

### Analyzer (`apps/analyzer`)

| Kategori       | Teknologi |
|----------------|-----------|
| Language       | Python 3.11 |
| Framework      | [FastAPI](https://fastapi.tiangolo.com/) |
| Pitch Extraction | [CREPE](https://github.com/marl/crepe) (tiny model) |
| Audio Processing | [librosa](https://librosa.org/) |
| YouTube Download | [yt-dlp](https://github.com/yt-dlp/yt-dlp) |

### Frontend (`apps/platform`)

| Kategori       | Teknologi |
|----------------|-----------|
| Language       | TypeScript 6.x |
| Framework      | React 19 |
| Build Tool     | [Vite 8](https://vite.dev/) |
| Routing        | [TanStack Router](https://tanstack.com/router) (file-based) |
| Data Fetching  | [TanStack Query](https://tanstack.com/query) |
| Styling        | [Tailwind CSS v4](https://tailwindcss.com/) |
| Icons          | [Lucide React](https://lucide.dev/) |
| API Codegen    | [Orval](https://orval.dev/) (Swagger → React Query hooks) |
| Linter         | [Biome](https://biomejs.dev/) |

### Monorepo Tooling

| Kategori       | Teknologi |
|----------------|-----------|
| Package Manager| [pnpm](https://pnpm.io/) (workspaces) |
| Build System   | [Nx](https://nx.dev/) |
| Containerization | Docker + Docker Compose |

---

## Prerequisites

- **Node.js** ≥ 20
- **pnpm** ≥ 9
- **Go** ≥ 1.25
- **Python** ≥ 3.11
- **PostgreSQL** ≥ 15
- **Docker** & Docker Compose (opsional, untuk deployment)
- **ffmpeg** (dibutuhkan oleh yt-dlp dan librosa)

---

## Getting Started

### 1. Install Dependencies

```bash
pnpm install
```

### 2. Setup Environment

```bash
cp .env.example .env
# Edit .env → isi ANTHROPIC_API_KEY
```

### 3. Setup Database

```bash
createdb maqam_detector_db
```

> Tabel akan otomatis di-migrate oleh GORM saat API pertama kali dijalankan.

### 4. Run Development

```bash
# Semua apps (API + Platform)
pnpm dev

# Masing-masing:
pnpm --filter api dev        # API (port 8080)
pnpm --filter platform dev   # Frontend (port 3000)

# Analyzer (Python):
cd apps/analyzer && uvicorn app.main:app --reload --port 8000
```

### 5. Docker (Production)

```bash
docker-compose up -d
```

---

## API Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/api/v1/analyze/youtube` | Analisis dari YouTube URL |
| `POST` | `/api/v1/analyze/upload` | Analisis dari file audio |
| `POST` | `/api/v1/analyze/record` | Analisis dari rekaman browser |
| `GET` | `/api/v1/analyses/:id` | Status & hasil analisis (polling) |
| `GET` | `/api/v1/history` | Riwayat analisis |
| `DELETE` | `/api/v1/history/:id` | Hapus riwayat |
| `GET` | `/api/v1/maqamat` | Daftar maqam yang didukung |
| `GET` | `/api/v1/maqamat/:id` | Detail maqam |

### Swagger Documentation

```
http://localhost:8080/swagger/index.html
```

---

## Available Scripts

| Command | Deskripsi |
|---------|-----------|
| `pnpm dev` | Jalankan semua apps (dev mode) |
| `pnpm build` | Build semua apps (production) |
| `pnpm --filter platform generate:api` | Generate API hooks dari Swagger |

---

## Documentation

Dokumentasi lengkap tersedia di `docs/core/`:

| File | Deskripsi |
|------|-----------|
| [PRD](docs/core/prd.md) | Product Requirements Document |
| [ERD](docs/core/erd.html) | Entity Relationship Diagram |
| [API Contract](docs/core/api-contract.md) | Spesifikasi API endpoint |
| [UI/UX Spec](docs/core/uiux-spec.md) | Design tokens, komponen, screen specs |
| [Implementation](docs/core/implementation.md) | Breakdown implementasi per fase |

---

## License

ISC — Anas (Cypress Consulting)

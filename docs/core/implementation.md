# Implementation Breakdown — MaqamDetector
**Versi:** 1.0.0  
**Tanggal:** Mei 2026  
**Scope:** Core Features F-001 s/d F-007  
**Base:** monorepo_gots_starterkit

---

## Prinsip Urutan Implementasi

```
Infrastructure → Data Layer → Analyzer Service → API → Frontend → Integration
```

Alasan urutan ini:
1. **Analyzer Service** adalah inti teknis paling kritis dan paling tidak pasti — dikerjakan paling awal agar risiko teknis terbuka lebih cepat
2. **API** bergantung pada Analyzer dan Database
3. **Frontend** bisa dikerjakan paralel setelah API contract tersedia (Swagger sudah ada sebagai kontrak)
4. **Integration** adalah validasi akhir seluruh pipeline

---

## Overview Fase

| Fase | Nama | Estimasi |
|------|------|----------|
| 0 | Project Setup | 1–2 hari |
| 1 | Python Analyzer Service | 4–6 hari |
| 2 | Database & Golang API | 4–5 hari |
| 3 | Frontend | 5–7 hari |
| 4 | Integration & Deployment | 2–3 hari |

**Total estimasi:** 16–23 hari kerja

---

## Fase 0 — Project Setup

### 0.1 Inisialisasi Monorepo

- [ ] Fork / clone `monorepo_gots_starterkit` sebagai base
- [ ] Rename project: `maqam-detector`
- [ ] Update `package.json` root: name, description, version
- [ ] Update `nx.json`: project name
- [ ] Rename `apps/platform` → tetap `platform` (atau sesuaikan preferensi)

### 0.2 Environment & Secrets

- [ ] Buat `.env` dari `.env.example` dengan variabel berikut:

```env
# API (Golang)
PORT=8080
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=
DB_NAME=maqam_detector_db
DB_PORT=5432
SSL_MODE=disable

# Analyzer (Python)
ANALYZER_PORT=8000
ANALYZER_BASE_URL=http://localhost:8000

# Claude API
ANTHROPIC_API_KEY=

# Frontend
VITE_API_URL=http://localhost:8080/api
```

- [ ] Tambahkan semua secret ke `.gitignore`
- [ ] Setup Husky pre-commit (Biome lint + format check)

### 0.3 Docker Compose

- [ ] Buat `docker-compose.yml` di root dengan service:

```yaml
services:
  postgres:
    image: postgres:15
    ports: ["5432:5432"]
    environment: { POSTGRES_DB, POSTGRES_USER, POSTGRES_PASSWORD }
    volumes: [postgres_data:/var/lib/postgresql/data]

  analyzer:
    build: ./apps/analyzer
    ports: ["8000:8000"]
    environment: { ANTHROPIC_API_KEY }

  api:
    build: ./apps/api
    ports: ["8080:8080"]
    depends_on: [postgres, analyzer]
    environment: { semua env API }

  platform:
    build: ./apps/platform
    ports: ["3000:3000"]
    depends_on: [api]
```

### 0.4 Verifikasi

- [ ] `pnpm dev` — pastikan `api` dan `platform` berjalan
- [ ] Swagger UI accessible di `http://localhost:8080/swagger/index.html`
- [ ] PostgreSQL terkoneksi

---

## Fase 1 — Python Analyzer Service

> **Mengapa dikerjakan duluan:** Ini adalah inti teknis paling tidak pasti. CREPE, yt-dlp, dan maqam matching harus divalidasi sebelum membangun API di atasnya.

### 1.1 Setup FastAPI

- [ ] Buat `apps/analyzer/` dengan struktur:

```
apps/analyzer/
├── app/
│   ├── main.py
│   ├── routes/
│   │   └── analyze.py
│   ├── services/
│   │   ├── pitch_extractor.py
│   │   ├── maqam_matcher.py
│   │   └── youtube_fetcher.py
│   └── data/
│       └── maqam_templates.json
├── requirements.txt
├── Dockerfile
└── .env
```

- [ ] Install dependencies:

```txt
# requirements.txt
fastapi==0.111.0
uvicorn==0.29.0
crepe==0.0.15
librosa==0.10.2
numpy==1.26.4
yt-dlp==2024.5.1
python-multipart==0.0.9
pydub==0.25.1
httpx==0.27.0
```

- [ ] Buat `Dockerfile` dengan base image `python:3.11-slim`:
  - Install `ffmpeg` (dependensi yt-dlp dan librosa)
  - Install requirements
  - Expose port 8000

### 1.2 Maqam Template Data

- [ ] Buat `maqam_templates.json` dengan 8 maqam core:

```json
[
  {
    "id": "hijaz",
    "name_latin": "Hijaz",
    "name_arabic": "حجاز",
    "tonic_semitone": 2,
    "intervals_semitones": [0, 1, 5, 7, 8, 11, 12],
    "pitch_class_profile": [0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0]
  }
]
```

  > Profil PCP (Pitch Class Profile) 12-bin untuk setiap maqam dari tonik D. Perlu dikerjakan dengan referensi teori maqam (MaqamWorld.com).

- [ ] Buat template untuk semua 8 maqam: Hijaz, Rast, Bayati, Nahawand, Kurd, Saba, Ajam, Jiharkah
- [ ] Validasi template dengan contoh audio yang sudah diketahui maqam-nya

### 1.3 YouTube Audio Fetcher

- [ ] Buat `youtube_fetcher.py`:

```python
# Fungsi: download audio-only dari YouTube URL
# - Validasi URL
# - Download ke temp file (UUID-named, /tmp/analyzer/)
# - Return path file + metadata (durasi, judul)
# - Cleanup otomatis setelah dipakai (finally block)
```

- [ ] Handle error: video private, geo-restricted, durasi terlalu panjang (>15 menit)
- [ ] Pastikan file temp selalu dihapus meskipun terjadi error

### 1.4 Pitch Extractor

- [ ] Buat `pitch_extractor.py`:

```python
# Fungsi: extract pitch dari file audio
# - Load audio dengan librosa (mono, sr=22050)
# - Crop segmen (segment_start, segment_duration)
# - Jalankan CREPE tiny model
# - Filter frame dengan confidence > 0.5
# - Konversi Hz → pitch class (0–11) dengan normalisasi 12-TET
# - Hitung Pitch Class Profile (PCP): histogram 12 bin, normalized
# - Return PCP array
```

- [ ] Test dengan file audio yang sudah diketahui maqam-nya
- [ ] Ukur waktu eksekusi (target < 15 detik untuk 60 detik audio)
- [ ] Validasi penggunaan RAM CREPE tiny (target < 400MB)

### 1.5 Maqam Matcher

- [ ] Buat `maqam_matcher.py`:

```python
# Fungsi: cocokkan PCP vs semua template
# - Load maqam_templates.json
# - Untuk setiap maqam: hitung cosine similarity antara input PCP dan template PCP
# - Coba semua 12 rotasi (transposisi) → ambil max similarity
# - Urutkan hasil (descending)
# - Return top 3: [{maqam_id, confidence_score, rank}]
```

- [ ] Implementasi rotasi 12 semitone untuk handle transposisi (lagu yang dimainkan di key berbeda)
- [ ] Normalisasi confidence score ke range 0.0–1.0

### 1.6 FastAPI Endpoint

- [ ] Buat `POST /internal/analyze` (hanya dipanggil oleh Golang API, tidak publik):

```python
# Request: { source_type, file_path, segment_start, segment_duration, mode }
# Response: { top3_candidates, audio_quality }
# - audio_quality: poor/fair/good/excellent (dari rata-rata CREPE confidence)
```

- [ ] Buat `GET /health` untuk health check

### 1.7 Validasi Akhir Fase 1

- [ ] Test manual: YouTube URL lagu Hijaz yang sudah diketahui → hasil harus Hijaz
- [ ] Test manual: Upload MP3 Bayati → hasil harus Bayati
- [ ] Test manual: URL dengan video private → error yang benar
- [ ] Ukur RAM usage saat CREPE berjalan (tidak boleh > 500MB)

---

## Fase 2 — Database & Golang API

### 2.1 GORM Models

- [ ] Buat `apps/api/model/session.go`:

```go
type Session struct {
    ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    CreatedAt    time.Time
    LastActiveAt time.Time
}
```

- [ ] Buat `apps/api/model/maqam.go`:

```go
type Maqam struct {
    ID                  string         `gorm:"primaryKey"` // slug: hijaz
    NameArabic          string
    NameLatin           string
    NameIndonesia       string
    PitchClassProfile   datatypes.JSON // []float64, 12 bins
    EmotionTags         datatypes.JSON // []string
    ExampleSongs        datatypes.JSON // []string
    IntervalDescription string
    TipsAransemen       string
    UpdatedAt           time.Time
}
```

- [ ] Buat `apps/api/model/analysis.go`:

```go
type Analysis struct {
    ID               uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    SessionID        uuid.UUID      `gorm:"type:uuid;not null"`
    Session          Session
    InputType        string         // youtube|upload|microphone|humming
    InputSource      string
    DetectedMaqamID  *string
    DetectedMaqam    *Maqam
    ConfidenceScore  *float64
    ExplanationText  string
    Status           string         // pending|processing|completed|failed
    ProcessingMs     *int
    ErrorCode        string
    ErrorMessage     string
    CreatedAt        time.Time
    CompletedAt      *time.Time
}
```

- [ ] Buat `apps/api/model/analysis_candidate.go`:

```go
type AnalysisCandidate struct {
    ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    AnalysisID      uuid.UUID `gorm:"type:uuid;not null"`
    MaqamID         string
    Maqam           Maqam
    ConfidenceScore float64
    Rank            int // 1, 2, 3
}
```

### 2.2 Database Migration & Seeder

- [ ] Konfigurasi GORM AutoMigrate di `config/database.go` untuk semua model
- [ ] Buat `seeders/maqam_seeder.go`: seed 8 maqam dari JSON file ke tabel `maqamat`
- [ ] Jalankan seeder saat `main.go` startup (hanya jika tabel kosong)
- [ ] Verifikasi: `SELECT * FROM maqamat;` mengembalikan 8 rows

### 2.3 Repository Layer

- [ ] `repository/session_repository.go`:
  - `FindOrCreate(sessionID uuid.UUID) (*Session, error)`
  - `UpdateLastActive(sessionID uuid.UUID) error`

- [ ] `repository/analysis_repository.go`:
  - `Create(analysis *Analysis) error`
  - `FindByID(id uuid.UUID) (*Analysis, error)`
  - `FindByIDAndSession(id, sessionID uuid.UUID) (*Analysis, error)`
  - `UpdateStatus(id uuid.UUID, status string) error`
  - `UpdateCompleted(id uuid.UUID, result *AnalysisResult) error`
  - `FindBySession(sessionID uuid.UUID, filter HistoryFilter) ([]Analysis, int64, error)`
  - `DeleteByIDAndSession(id, sessionID uuid.UUID) error`

- [ ] `repository/maqam_repository.go`:
  - `FindAll() ([]Maqam, error)`
  - `FindByID(id string) (*Maqam, error)`

### 2.4 Middleware

- [ ] `middleware/cors.go`: Allow origins dari `VITE_API_URL`, methods, headers termasuk `X-Session-ID`
- [ ] `middleware/session.go`: Validasi `X-Session-ID` header (harus UUID valid), inject ke context
- [ ] `middleware/rate_limiter.go`: In-memory rate limit per session ID (10 analisis/jam, 5 upload/jam)
- [ ] `middleware/request_logger.go`: Sudah ada dari starter kit (Logrus) — sesuaikan format

### 2.5 Service: Analyze

- [ ] `service/analyze_service.go`:

```
AnalyzeYoutube(sessionID, url, segmentStart, segmentDuration):
  1. Buat record Analysis (status: pending)
  2. Jalankan goroutine: processAnalysis()
  3. Return analysis_id + estimated_seconds

processAnalysis(analysisID):
  1. Update status → processing
  2. POST ke Python Analyzer /internal/analyze
  3. Ambil top3 candidates dari response
  4. POST ke Claude API → generate explanation (Bahasa Indonesia)
  5. Simpan hasil ke DB (candidates + explanation)
  6. Update status → completed
  7. Jika error di langkah manapun → update status → failed + error_code
```

- [ ] `service/claude_service.go`:
  - Build prompt dari top1 maqam + top3 candidates
  - Call Anthropic API (model: claude-sonnet-4-6)
  - Parse response → struct explanation
  - Handle API error + timeout (max 30 detik)

- [ ] `service/history_service.go`:
  - `GetHistory(sessionID, filter)` → paginated list
  - `DeleteHistory(sessionID, analysisID)` → validasi ownership

- [ ] `service/maqam_service.go`:
  - `GetAll()` → list semua maqam
  - `GetByID(id)` → detail satu maqam

### 2.6 DTO & Swagger Annotations

- [ ] Buat semua DTO request/response di `dto/`:
  - `AnalyzeYoutubeRequest`, `AnalyzeUploadRequest`, `AnalyzeRecordRequest`
  - `AnalyzeResponse` (202)
  - `AnalysisDetailResponse` (hasil polling)
  - `HistoryListResponse`, `HistoryItemResponse`
  - `MaqamListResponse`, `MaqamDetailResponse`
  - `SuccessResponse`, `ErrorResponse` (sudah ada di starter kit)

- [ ] Tambahkan Swagger annotations di semua DTO dan handler
- [ ] Generate docs: `swag init -g main.go`
- [ ] Verifikasi Swagger UI menampilkan semua 9 endpoint dengan benar

### 2.7 Handler Layer

- [ ] `handler/analyze_handler.go`:
  - `POST /api/v1/analyze/youtube`
  - `POST /api/v1/analyze/upload` (multipart)
  - `POST /api/v1/analyze/record` (multipart)
  - `GET /api/v1/analyses/:id`

- [ ] `handler/history_handler.go`:
  - `GET /api/v1/history`
  - `DELETE /api/v1/history/:id`

- [ ] `handler/maqam_handler.go`:
  - `GET /api/v1/maqamat`
  - `GET /api/v1/maqamat/:id`

### 2.8 Validasi Akhir Fase 2

- [ ] Test semua endpoint via Swagger UI
- [ ] `POST /analyze/youtube` → 202, analysis_id kembali
- [ ] `GET /analyses/:id` → status berubah dari pending → processing → completed
- [ ] `GET /history` → data muncul
- [ ] `DELETE /history/:id` → data terhapus
- [ ] `GET /maqamat` → 8 maqam
- [ ] Rate limit bekerja setelah 10 request/jam

---

## Fase 3 — Frontend

### 3.1 Setup & Konfigurasi

- [ ] Verifikasi Orval config di `apps/platform/orval.config.ts`:

```ts
export default defineConfig({
  maqam: {
    input: { target: 'http://localhost:8080/swagger/doc.json' },
    output: {
      target: 'src/api/endpoints',
      schemas: 'src/api/model',
      client: 'react-query',
      mode: 'tags-split',
    },
  },
})
```

- [ ] Generate API hooks: `pnpm --filter platform generate:api`
- [ ] Verifikasi hooks ter-generate: `useGetAnalysesAnalysisId`, `usePostAnalyzeYoutube`, dll

### 3.2 Design Tokens → Tailwind Config

- [ ] Update `tailwind.config.ts` dengan semua color tokens dari UX/UI Spec:

```ts
theme: {
  extend: {
    colors: {
      brand: {
        primary: '#0F6E56',
        'primary-hover': '#085041',
        'primary-subtle': '#E1F5EE',
        secondary: '#534AB7',
        'secondary-subtle': '#EEEDFE',
      },
      confidence: {
        'very-high': '#3B6D11',
        high: '#639922',
        medium: '#854F0B',
        low: '#993C1D',
        'very-low': '#A32D2D',
      }
    },
    fontFamily: {
      arabic: ['Amiri', 'serif'],
    },
    borderRadius: {
      xl: '16px',
      '2xl': '24px',
    }
  }
}
```

- [ ] Import font Amiri di `index.html` (Google Fonts)
- [ ] Setup dark mode: `darkMode: 'class'`

### 3.3 Komponen Atoms

Buat di `src/components/atoms/`:

- [ ] `Button.tsx` — variant: filled/outline/ghost/text · size: sm/md/lg · state: default/loading/disabled
- [ ] `Input.tsx` — state: default/focus/valid/error/disabled
- [ ] `Badge.tsx` — variant: brand/success/warning/error/info/neutral
- [ ] `Chip.tsx` — variant: outline/filled
- [ ] `Spinner.tsx` — size: sm/md/lg · color: brand/white/neutral
- [ ] `ProgressBar.tsx` — variant: brand/success/warning/error · animated fill
- [ ] `SkeletonLoader.tsx` — variant: text/image/card

### 3.4 Komponen Molecules

Buat di `src/components/molecules/`:

- [ ] `Toast.tsx` + `useToast.ts` — variant: success/error/info/warning · auto-dismiss
- [ ] `ConfidenceBadge.tsx` — tampil label + warna sesuai nilai confidence
- [ ] `EmptyState.tsx` — ikon + judul + sub + optional CTA
- [ ] `ContextMenu.tsx` — trigger + dropdown item list

### 3.5 Feature: Analyzer (Home Screen)

Buat di `src/features/analyzer/`:

#### Input Method Tabs

- [ ] `InputMethodTabs.tsx` — tab switcher YouTube/Upload/Rekam
- [ ] State management: tab aktif, form data per tab (gunakan `useState` lokal atau Jotai jika perlu share antar komponen)

#### Tab YouTube

- [ ] `YoutubeInputPanel.tsx`:
  - URL input dengan validasi real-time (debounce 800ms)
  - `VideoPreviewCard.tsx` — thumbnail + judul + durasi (fetch dari YouTube oEmbed API)
  - Segment selector (start + duration)

#### Tab Upload

- [ ] `FileUploadZone.tsx`:
  - Drag & drop handler
  - File validation (MIME type, ukuran ≤ 50MB)
  - `FileInfoCard.tsx` — tampil nama file + ukuran + durasi

#### Tab Rekam

- [ ] `RecordPanel.tsx`:
  - Mode toggle (Mikrofon/Humming)
  - `WaveformVisualizer.tsx` — WebAudio API, real-time bar animation
  - `RecordButton.tsx` — state: idle/recording/done dengan timer counter
  - Durasi validator (minimum 5 detik)

#### Analyze Button & Submit Logic

- [ ] `AnalyzeButton.tsx` — disabled/loading state
- [ ] Submit handler:
  - YouTube: `mutate` dari `usePostAnalyzeYoutube`
  - Upload: `mutate` dari `usePostAnalyzeUpload` (FormData)
  - Record: `mutate` dari `usePostAnalyzeRecord` (FormData + mode)
  - Setelah 202 → navigate ke `/processing/:analysisId`

#### Supported Maqam Row

- [ ] `SupportedMaqamRow.tsx` — fetch dari `useGetMaqamat`, tampil sebagai chip row

### 3.6 Feature: Processing Screen

Buat di `src/features/analyzer/` (atau `src/routes/processing.$analysisId.tsx`):

- [ ] `ProcessingScreen.tsx`:
  - Tampil info sumber audio (thumbnail/ikon + judul)
  - `AnalysisProgressSteps.tsx` — 4 step dengan state pending/active/completed
  - Estimasi waktu tersisa
  - Tombol Batalkan

- [ ] **Polling logic** dengan TanStack Query:

```ts
const { data } = useGetAnalysesAnalysisId(analysisId, {
  refetchInterval: (data) =>
    data?.data?.status === 'pending' || data?.data?.status === 'processing'
      ? 2000
      : false,
  refetchIntervalInBackground: false,
})

// Timeout 90 detik → navigate ke error
```

- [ ] Auto-navigate ke `/result/:analysisId` jika status `completed`
- [ ] Auto-navigate ke `/error` jika status `failed`

### 3.7 Feature: Hasil Analisis (Result Screen)

Buat di `src/features/result/`:

- [ ] `ResultScreen.tsx` — layout 2 kolom (desktop), 1 kolom (mobile)
- [ ] `MaqamResultCard.tsx`:
  - Nama Latin (large) + nama Arab (Amiri font)
  - Confidence progress bar dengan warna dinamis
  - Label verbal confidence
  - Interval tangga nada (monospace)

- [ ] `MaqamCandidateList.tsx`:
  - Tampil rank 2 dan 3
  - `ConfidenceWarningBlock.tsx` — muncul jika confidence < 70%

- [ ] `MaqamExplanation.tsx` — accordion (mobile), expanded (desktop):
  - `CharacteristicsSection.tsx` — teks + emotion tags
  - `ScaleStructureSection.tsx` — interval monospace + penjelasan
  - `ExampleSongsSection.tsx` — bullet list
  - `BanjariTipsSection.tsx` — highlighted block

- [ ] Action button: "Analisis Lagu Lain" → navigate ke `/`

### 3.8 Feature: Riwayat Analisis (History Screen)

Buat di `src/features/history/`:

- [ ] `HistoryScreen.tsx`:
  - `FilterBar.tsx` — Semua/YouTube/Upload/Rekaman
  - List `HistoryItem.tsx`
  - Pagination atau infinite scroll (TanStack Query `useInfiniteQuery`)
  - `EmptyState` jika kosong

- [ ] `HistoryItem.tsx`:
  - Thumbnail/ikon input type
  - Judul sumber + maqam + confidence + tanggal
  - `ContextMenu` dengan opsi: Lihat Detail, Hapus
  - Status badge jika bukan completed

- [ ] Hapus item: optimistic update dengan `useMutation` + invalidate query

### 3.9 Error States & Toast

- [ ] `ErrorScreen.tsx` — tampil jika analysis failed + tombol Coba Lagi/Kembali
- [ ] Error message mapping (dari API error code → pesan Bahasa Indonesia)
- [ ] Global `ToastProvider` di root layout
- [ ] Trigger toast: setelah hapus riwayat, setelah error kecil

### 3.10 Routing Setup

- [ ] Konfigurasi file-based routes di `src/routes/`:

```
routes/
├── __root.tsx           ← Root layout (Navbar, Toast provider)
├── index.tsx            ← Home / Analyzer
├── processing.$analysisId.tsx  ← Processing screen
├── result.$analysisId.tsx      ← Hasil analisis
└── history.tsx          ← Riwayat
```

- [ ] Setup `X-Session-ID` di Axios custom instance (`src/api/mutator/custom-instance.ts`):

```ts
// Generate UUID sekali, simpan di localStorage
// Inject sebagai header di setiap request
const sessionId = localStorage.getItem('session_id') ?? crypto.randomUUID()
localStorage.setItem('session_id', sessionId)
instance.defaults.headers['X-Session-ID'] = sessionId
```

### 3.11 Validasi Akhir Fase 3

- [ ] Home screen — 3 tab berfungsi, semua input valid/error state tampil
- [ ] Submit YouTube → navigate ke Processing
- [ ] Processing — progress steps update, auto-navigate ke Result setelah selesai
- [ ] Result — semua section tampil, confidence color benar
- [ ] History — list tampil, filter bekerja, delete bekerja
- [ ] Responsive: tampilan mobile (390px) dan desktop (1440px) sesuai spec
- [ ] Dark mode bekerja (jika diimplementasikan)

---

## Fase 4 — Integration & Deployment

### 4.1 End-to-End Test Manual

Jalankan semua happy path secara berurutan:

- [ ] YouTube → lagu Hijaz yang diketahui → hasil Hijaz, confidence > 70%
- [ ] YouTube → lagu Bayati → hasil Bayati
- [ ] Upload MP3 → hasil sesuai
- [ ] Rekam mikrofon (nyanyikan melodi Rast) → hasil mendekati Rast
- [ ] History → semua analisis sebelumnya muncul
- [ ] Delete history → item hilang
- [ ] Rate limit → setelah 10 analisis, error 429 + pesan yang benar
- [ ] Error case: YouTube private → error message yang benar
- [ ] Error case: file > 50MB → error message yang benar

### 4.2 VPS Deployment

#### Persiapan Server

- [ ] Install Docker + Docker Compose di VPS
- [ ] Install Nginx sebagai reverse proxy
- [ ] Install Certbot (Let's Encrypt SSL)
- [ ] Setup domain/subdomain (misal: `maqam.example.com` dan `api.maqam.example.com`)

#### Nginx Config

```nginx
# Frontend (SPA)
server {
    server_name maqam.example.com;
    root /var/www/platform;
    try_files $uri $uri/ /index.html;  # SPA fallback
}

# API
server {
    server_name api.maqam.example.com;
    location / { proxy_pass http://localhost:8080; }
}
```

#### Deploy Steps

- [ ] Build frontend: `pnpm --filter platform build` → upload `dist/` ke `/var/www/platform`
- [ ] Build dan jalankan Docker Compose: `docker-compose up -d`
- [ ] Setup SSL: `certbot --nginx -d maqam.example.com -d api.maqam.example.com`
- [ ] Verifikasi semua service berjalan: `docker-compose ps`
- [ ] Cek RAM usage: `free -h` (harus ada headroom ≥ 200MB)

#### Post-Deploy Checklist

- [ ] Akses `https://maqam.example.com` → Home screen tampil
- [ ] Test analisis YouTube dari browser production
- [ ] Cek log error: `docker-compose logs -f api`
- [ ] Cek log analyzer: `docker-compose logs -f analyzer`
- [ ] Verifikasi temp file audio terhapus setelah analisis selesai

### 4.3 Monitoring Minimal (Fase Internal)

- [ ] Setup log rotation untuk Docker logs
- [ ] Cek VPS resource usage setelah 1 hari penggunaan normal

---

## Dependency Map Antar Fase

```
Fase 0 (Setup)
    │
    ├──→ Fase 1 (Python Analyzer)  ← bisa mulai langsung
    │         │
    │         ▼
    └──→ Fase 2 (Golang API)  ← butuh Analyzer selesai
              │
              ├──→ Generate Swagger → Orval
              │
              ▼
         Fase 3 (Frontend)  ← bisa paralel setelah Fase 2 sebagian
              │
              ▼
         Fase 4 (Integration & Deploy)
```

**Yang bisa dikerjakan paralel:**
- Fase 1 dan Fase 0 → bisa bersamaan
- Fase 3.1–3.4 (setup + atoms + molecules) → bisa mulai setelah Swagger tersedia, sebelum Fase 2 selesai penuh

---

## Risiko & Mitigasi

| Risiko | Probabilitas | Dampak | Mitigasi |
|--------|-------------|--------|---------|
| CREPE OOM di VPS 2GB RAM | Sedang | Tinggi | Gunakan model `tiny`, test RAM sebelum deploy |
| yt-dlp blocked oleh YouTube | Rendah | Tinggi | Cache hasil analisis, hindari re-download URL yang sama |
| Akurasi maqam rendah (< 60%) | Sedang | Sedang | Mulai dengan maqam yang paling distinctive (Hijaz), iterate template |
| Claude API timeout | Rendah | Rendah | Timeout 30 detik + fallback teks statis |
| CORS issue setelah deploy | Rendah | Sedang | Test CORS config di staging sebelum production |

---

## Catatan Penting

### Tentang Akurasi Maqam

Akurasi deteksi sangat bergantung pada kualitas **pitch class profile template** di `maqam_templates.json`. Tahap paling kritis adalah **validasi template** di Fase 1.2 — disarankan:

1. Kumpulkan minimal 3–5 lagu per maqam yang sudah pasti maqam-nya
2. Ekstrak PCP dari setiap lagu
3. Rata-ratakan → jadikan template
4. Iterasi jika akurasi kurang dari 70%

### Tentang Claude Explanation

Prompt ke Claude perlu di-craft dengan baik agar output konsisten. Sertakan:
- Nama maqam yang terdeteksi
- Confidence score
- Top 3 kandidat
- Instruksi: Bahasa Indonesia, sertakan tips aransemen banjari

### File Cleanup Policy

Setiap temp file audio **wajib** dihapus di `finally` block Python, baik analisis berhasil maupun gagal. Ini kritis karena storage VPS hanya 20GB.

---

**Versi History:**

| Versi | Tanggal | Perubahan |
|-------|---------|-----------|
| 1.0.0 | Mei 2026 | Draft awal — core features F-001 s/d F-007 |

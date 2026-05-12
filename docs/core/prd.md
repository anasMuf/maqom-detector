# Product Requirements Document
# MaqamDetector — Pendeteksi Maqam Musik Arab & Timur Tengah

**Versi:** 1.0.0  
**Tanggal:** Mei 2026  
**Status:** Draft  
**Author:** Anas (Cypress Consulting)

---

## Daftar Isi

1. [Product Overview](#1-product-overview)
2. [Background & Problem Statement](#2-background--problem-statement)
3. [Solution Overview](#3-solution-overview)
4. [Goals & Success Metrics](#4-goals--success-metrics)
5. [Features](#5-features)
6. [User Stories](#6-user-stories)
7. [Technical Architecture](#7-technical-architecture)
8. [Non-Functional Requirements](#8-non-functional-requirements)
9. [Assumptions & Constraints](#9-assumptions--constraints)
10. [Open Questions](#10-open-questions)
11. [References](#11-references)

---

## 1. Product Overview

**MaqamDetector** adalah aplikasi web yang membantu musisi, arranger, dan penggemar musik banjari untuk mengidentifikasi maqam (tangga nada modal) dari lagu-lagu Arab dan Timur Tengah secara otomatis — baik dari file audio, link YouTube, rekaman mikrofon, maupun humming melodi.

Output aplikasi ini berupa nama maqam yang terdeteksi beserta penjelasan karakteristiknya dalam Bahasa Indonesia, dilengkapi confidence score dan kandidat alternatif. Hal ini membantu komunitas banjari dalam proses eksplorasi dan aransemen lagu-lagu referensi dari munsyid maupun grup dari kawasan Syiria, Mesir, dan Timur Tengah lainnya.

---

## 2. Background & Problem Statement

### 2.1 Konteks Komunitas

Komunitas banjari (musik rebana/hadrah) di Indonesia, terutama yang aktif dalam kompetisi antar grup, kerap mengadopsi lagu-lagu dari munsyid dan grup musik Arab Timur Tengah — terutama dari Syiria — sebagai referensi aransemen. Proses ini melibatkan:

- Mendengarkan lagu referensi dari YouTube atau platform streaming
- Memahami struktur melodi dan tangga nada
- Melakukan aransemen ulang sesuai instrumen banjari

### 2.2 Problem Statement

Dalam proses eksplorasi referensi tersebut, anggota komunitas sering kali menghadapi kebingungan mendasar:

> *"Ini maqam apa ya? Hijaz atau Rast?"*

Kebingungan ini timbul karena:

1. **Tidak ada tools otomatis** — Belum ada aplikasi consumer-ready yang bisa mengidentifikasi maqam dari audio secara langsung. Tools seperti Shazam atau AHA Music hanya mengenali judul lagu, bukan maqam-nya.

2. **Pengetahuan teori maqam tidak merata** — Tidak semua anggota komunitas memiliki latar belakang teori musik Arab yang memadai untuk mengidentifikasi maqam secara auditori.

3. **Maqam bersifat mikrotonal** — Sistem maqam menggunakan interval quarter-tone yang tidak ada dalam notasi musik Barat standar, sehingga referensi teori Barat tidak cukup untuk membantu.

4. **Proses manual membutuhkan pakar** — Satu-satunya cara saat ini adalah bertanya kepada musisi senior yang memahami teori maqam, yang tidak selalu tersedia.

5. **Kebutuhan mendesak sebelum lomba** — Grup yang sedang mempersiapkan aransemen untuk kompetisi membutuhkan jawaban cepat.

### 2.3 Dampak Problem

- Proses aransemen menjadi lebih lambat dan tidak efisien
- Hasil aransemen bisa "salah maqam" jika identifikasi manual keliru
- Anggota baru komunitas sulit berkontribusi dalam eksplorasi repertoar

---

## 3. Solution Overview

MaqamDetector menggabungkan tiga lapisan teknologi:

```
Audio Input
    ↓
Pitch Extraction (CREPE / librosa)
    ↓
Maqam Template Matching (Rule-based, 24-TET aware)
    ↓
AI Explanation (Claude API)
    ↓
Output: Nama Maqam + Penjelasan + Confidence Score
```

**Pendekatan utama:** Rule-based template matching berbasis pitch class profile (PCP) yang dipetakan ke 24-TET (quarter-tone aware), dikombinasikan dengan Claude API untuk menghasilkan penjelasan kontekstual dalam Bahasa Indonesia. Pendekatan ini dipilih karena:

- Tidak membutuhkan dataset audio berlabel yang besar
- Transparan dan dapat di-debug (bukan black-box ML)
- Mudah di-extend dengan menambahkan template maqam baru
- Lebih efisien dari sisi komputasi untuk VPS dengan resource terbatas

**Maqam yang didukung (versi awal):** Hijaz, Rast, Bayati, Nahawand, Kurd, Saba, Ajam, Jiharkah — 8 maqam paling umum dalam repertoar musik Arab Timur Tengah yang sering dijadikan referensi banjari.

### 3.1 Deteksi Vokal vs Instrumen

Sistem mendeteksi maqam dari **melodi dominan dalam audio** — baik bersumber dari vokal maupun instrumen melodi. CREPE sebagai engine pitch extraction memang awalnya dirancang untuk *vocal pitch estimation*, sehingga performa terbaik justru pada audio yang mengandung vokal sebagai pembawa melodi utama.

Ini sangat relevan untuk konteks komunitas banjari, mengingat banyak referensi berasal dari grup yang **tidak menggunakan instrumen melodi** (hanya vokal + rebana/perkusi). Dalam kondisi ini, karena rebana tidak memiliki pitch, CREPE dapat fokus sepenuhnya pada vokal sehingga hasil deteksi cenderung lebih bersih dan akurat.

| Skenario Audio Referensi | Kualitas Deteksi | Keterangan |
|--------------------------|-----------------|------------|
| Vokal + rebana/perkusi saja | ✅ Ideal | Tidak ada interferensi instrumen melodi; paling optimal |
| Vokal + instrumen melodi (oud, qanun, dll.) | ✅ Baik | CREPE mengambil dominant pitch, biasanya vokal lead |
| Instrumen melodi saja (tanpa vokal) | ✅ Baik | Cocok untuk referensi instrumental Syiria/Mesir |
| Multi-vokal / paduan suara | ⚠️ Cukup | Pitch yang diambil adalah yang paling dominan, belum tentu melodi utama |
| Audio kualitas rendah / noise tinggi | ⚠️ Kurang | Confidence score akan rendah; disarankan cari sumber audio yang lebih baik |

---

## 4. Goals & Success Metrics

### 4.1 Goals

| Tujuan | Deskripsi |
|--------|-----------|
| **Akurasi** | Deteksi maqam yang benar untuk 8 maqam utama |
| **Kemudahan** | Anggota komunitas tanpa latar belakang teori bisa menggunakan |
| **Kecepatan** | Hasil analisis dalam waktu < 30 detik untuk lagu berdurasi 3–5 menit |
| **Aksesibilitas** | Bisa diakses dari perangkat apapun (mobile-friendly) |

### 4.2 Success Metrics (Internal Phase)

- **Akurasi deteksi** ≥ 75% pada 8 maqam utama (divalidasi komunitas)
- **Waktu analisis** < 30 detik untuk input YouTube
- **User satisfaction** ≥ 4/5 berdasarkan feedback internal komunitas
- **Adoption rate** ≥ 60% anggota aktif komunitas menggunakan dalam 1 bulan pertama

### 4.3 Success Metrics (Public Phase)

- Monthly Active Users (MAU) ≥ 200 dalam 3 bulan pertama
- Bounce rate < 50%
- Return rate ≥ 40% (user kembali menggunakan dalam 7 hari)

---

## 5. Features

### 5.1 Core Features (Must Have — Phase 1)

#### F-001 · Analisis dari Link YouTube
Pengguna dapat menempelkan URL YouTube (video klip atau audio lagu Arab/Timur Tengah) dan mendapatkan deteksi maqam secara otomatis.

- Input: URL YouTube valid
- Proses: Ekstraksi audio via yt-dlp → pitch extraction → maqam matching
- Output: Nama maqam + confidence score + top 3 kandidat

#### F-002 · Analisis dari Upload File Audio
Pengguna dapat mengunggah file audio lokal (MP3, WAV, M4A, FLAC, OGG) untuk dianalisis.

- Batas ukuran file: maks. 50MB
- Durasi yang dianalisis: 60 detik pertama (configurable)

#### F-003 · Rekam Langsung dari Mikrofon Browser
Pengguna dapat merekam audio langsung dari browser (melodi yang dinyanyikan, dilantunkan, atau dimainkan secara langsung).

- Durasi rekaman: 10–60 detik
- Real-time waveform visualization saat merekam

#### F-004 · Humming / Senandung Melodi
Pengguna dapat senandung (humming) atau menyanyikan melodi dari lagu yang ingin diidentifikasi maqam-nya.

- Sama dengan F-003 secara teknis, namun dengan toleransi pitch yang lebih longgar
- Mode khusus "humming" dengan threshold similarity yang disesuaikan

#### F-005 · Output Penjelasan Maqam (Bahasa Indonesia)
Setelah deteksi, sistem menghasilkan penjelasan komprehensif dalam Bahasa Indonesia yang mencakup:

- Nama maqam (Arab + transliterasi)
- Karakteristik emosi dan nuansa
- Struktur interval/tangga nada
- Contoh lagu/munsyid yang terkenal menggunakan maqam ini
- Tips relevan untuk aransemen banjari

#### F-006 · Confidence Score & Kandidat Alternatif
Menampilkan tingkat keyakinan deteksi (dalam %) dan menampilkan 2–3 kandidat maqam alternatif jika confidence score maqam utama tidak terlalu tinggi (di bawah 70%).

#### F-007 · Riwayat Analisis
Menyimpan riwayat analisis yang sudah dilakukan per pengguna (session-based untuk guest, per akun untuk registered user).

- Tampil maks. 20 entri terakhir
- Bisa diulangi atau dihapus

---

### 5.2 Nice to Have (Phase 2)

#### F-008 · Visualisasi Pitch & Tangga Nada
Menampilkan visualisasi grafis pitch yang terdeteksi dari audio, dioverlay dengan struktur tangga nada maqam yang teridentifikasi.

#### F-009 · Perbandingan Dua Maqam
Fitur "compare mode" yang menampilkan dua maqam berdampingan — interval, emosi, karakteristik — untuk membantu arranger memilih antara dua kandidat maqam.

#### F-010 · Bookmark & Koleksi Pribadi
Pengguna dapat menyimpan hasil analisis ke dalam koleksi bertema (misalnya: "Referensi Lomba 2026", "Maqam Hijaz Favorites").

#### F-011 · Share Hasil Analisis
Pengguna dapat berbagi hasil analisis via link atau card gambar yang bisa dibagikan ke media sosial atau WhatsApp grup komunitas.

#### F-012 · Mode Belajar Maqam
Konten edukasi interaktif tentang karakteristik tiap maqam — dilengkapi contoh audio dan latihan identifikasi telinga (ear training).

#### F-013 · Koreksi Komunitas
Pengguna dapat menandai hasil deteksi sebagai "kurang tepat" dan memberikan koreksi. Data dikumpulkan sebagai referensi peningkatan akurasi template — bersifat self-moderated oleh komunitas, tanpa expert reviewer di fase ini. Admin dapat melihat log koreksi dan mempertimbangkannya untuk update template secara manual.

---

### 5.3 Plan to Have (Phase 3 — Public)

#### F-014 · Autentikasi & Profil Pengguna
Sistem login (email/password atau Google OAuth) dengan profil pengguna yang menyimpan preferensi, riwayat, dan koleksi.

#### F-015 · API Publik
REST API yang bisa diakses developer eksternal untuk integrasi maqam detection ke aplikasi pihak ketiga.

#### F-016 · Deteksi Multi-Maqam dalam Satu Lagu
Identifikasi perubahan maqam (modulasi) dalam satu lagu — misalnya lagu yang dimulai dari Rast lalu berpindah ke Hijaz.

#### F-017 · Ekspansi Database Maqam
Menambahkan maqam-maqam yang lebih jarang: Sikah, Huzam, Awj, Nawa Athar, Maqam Rast Syam, dll.

#### F-018 · Rekomendasi Lagu Banjari Serupa
Berdasarkan maqam yang terdeteksi, sistem merekomendasikan lagu banjari populer yang menggunakan maqam yang sama sebagai referensi aransemen.

---

## 6. User Stories

### Role Definitions

| Role | Deskripsi |
|------|-----------|
| **Anggota** | Anggota komunitas banjari, bisa jadi pemain atau pendengar umum |
| **Arranger** | Anggota yang bertugas menyusun aransemen lagu untuk lomba |
| **Admin** | Pengelola internal aplikasi (bisa dari komunitas atau developer) |
| **Guest** | Pengunjung publik tanpa akun (fase public) |

---

### F-001 · Analisis YouTube

**US-001**  
*Sebagai* **Arranger**,  
*Saya ingin* menempelkan link YouTube lagu Syiria yang sedang saya jadikan referensi,  
*Agar* saya bisa langsung tahu maqam apa yang digunakan tanpa harus menganalisis sendiri.

**Acceptance Criteria:**
- URL YouTube divalidasi sebelum diproses
- Durasi video maksimum 15 menit (untuk mencegah abuse dan menghemat resource)
- Proses analisis berjalan dalam background dengan loading indicator
- Jika video private/tidak tersedia, tampil pesan error yang informatif
- Hasil muncul dalam ≤ 30 detik

---

**US-002**  
*Sebagai* **Anggota**,  
*Saya ingin* melihat loading progress saat analisis YouTube sedang berjalan,  
*Agar* saya tahu sistem sedang bekerja dan tidak mengira aplikasi hang.

**Acceptance Criteria:**
- Progress bar atau step indicator yang menunjukkan tahap: Mengunduh audio → Menganalisis pitch → Mencocokkan maqam → Menghasilkan penjelasan
- Estimasi waktu tersisa ditampilkan jika memungkinkan

---

### F-002 · Upload File Audio

**US-003**  
*Sebagai* **Arranger**,  
*Saya ingin* mengunggah file MP3 lagu yang sudah saya download dari munsyid favorit saya,  
*Agar* saya bisa menganalisis maqam-nya meski tidak punya link YouTube.

**Acceptance Criteria:**
- Format yang didukung: MP3, WAV, M4A, FLAC, OGG
- Ukuran maksimum: 50MB
- File ditolak dengan pesan error jika format tidak sesuai atau ukuran melebihi batas
- File audio dihapus dari server setelah proses selesai (privacy)

---

### F-003 & F-004 · Mikrofon & Humming

**US-004**  
*Sebagai* **Anggota**,  
*Saya ingin* merekam saya senandung melodi yang sedang ada di kepala saya,  
*Agar* saya bisa tahu itu maqam apa meski saya tidak punya filenya.

**Acceptance Criteria:**
- Browser meminta izin mikrofon sebelum merekam
- Waveform audio tampil secara real-time saat merekam
- Ada tombol Start, Stop, dan Re-record
- Durasi rekaman minimum 5 detik, maksimum 60 detik
- Peringatan ditampilkan jika mikrofon terlalu sunyi atau ada noise berlebihan

---

**US-005**  
*Sebagai* **Anggota**,  
*Saya ingin* mode humming yang lebih toleran terhadap ketidaksempurnaan nada,  
*Agar* hasil deteksi tetap akurat meski senandung saya tidak selalu tepat secara pitch.

**Acceptance Criteria:**
- Mode "Humming" memiliki threshold similarity yang lebih longgar dibanding mode audio normal
- Sistem tetap memberikan top 3 kandidat maqam dengan penjelasan uncertainty-nya

---

### F-005 · Output Penjelasan

**US-006**  
*Sebagai* **Arranger**,  
*Saya ingin* membaca penjelasan karakteristik maqam yang terdeteksi dalam Bahasa Indonesia,  
*Agar* saya bisa memahami nuansa emosi dan karakter melodi maqam tersebut sebelum mulai aransemen.

**Acceptance Criteria:**
- Penjelasan mencakup: nama maqam, karakteristik emosi, struktur interval, contoh lagu terkenal
- Bahasa yang digunakan natural dan mudah dipahami (bukan akademis)
- Ada section khusus "Tips untuk Aransemen Banjari" yang relevan

---

**US-007**  
*Sebagai* **Anggota**,  
*Saya ingin* melihat struktur tangga nada maqam secara visual,  
*Agar* saya bisa memahami interval-intervalnya secara lebih intuitif.

**Acceptance Criteria:**
- Tampil representasi visual tangga nada (notasi sederhana atau diagram)
- Menandai note-note khas maqam (seperti half-flat pada Bayati)

---

### F-006 · Confidence Score

**US-008**  
*Sebagai* **Arranger**,  
*Saya ingin* melihat seberapa yakin sistem terhadap hasil deteksinya,  
*Agar* saya bisa memutuskan apakah perlu verifikasi lebih lanjut atau tidak.

**Acceptance Criteria:**
- Confidence score ditampilkan dalam persentase (%)
- Jika confidence < 70%, tampil peringatan bahwa hasil kurang pasti dan tampilkan kandidat alternatif
- Penjelasan singkat mengapa confidence rendah (misalnya: "melodi terlalu pendek" atau "karakteristik maqam ambigu")

---

### F-007 · Riwayat Analisis

**US-009**  
*Sebagai* **Anggota**,  
*Saya ingin* melihat riwayat analisis yang pernah saya lakukan,  
*Agar* saya tidak perlu menganalisis ulang lagu yang sama.

**Acceptance Criteria:**
- Tampil daftar 20 analisis terakhir dengan: judul/sumber, maqam terdeteksi, tanggal
- Bisa klik untuk melihat ulang hasil detail
- Bisa hapus entri individual dari riwayat

---

### F-013 · Koreksi Komunitas

**US-010**  
*Sebagai* **Anggota** yang memahami teori maqam,  
*Saya ingin* memberikan koreksi jika hasil deteksi salah,  
*Agar* sistem bisa terus berkembang dan data referensi komunitas semakin akurat.

**Acceptance Criteria:**
- Tombol "Kurang Tepat / Laporkan Koreksi" tersedia di setiap hasil
- User bisa memilih maqam yang menurut mereka benar dari dropdown
- Admin bisa mereview koreksi dan memvalidasi sebelum digunakan untuk improve sistem

---

### Admin

**US-011**  
*Sebagai* **Admin**,  
*Saya ingin* melihat statistik penggunaan aplikasi (jumlah analisis, maqam yang paling sering dideteksi, error rate),  
*Agar* saya bisa memantau kesehatan sistem dan tren penggunaan komunitas.

**Acceptance Criteria:**
- Dashboard admin dengan grafik usage harian/mingguan
- Tabel maqam terpopuler
- Log error dan waktu pemrosesan rata-rata

---

## 7. Technical Architecture

### 7.1 Rekomendasi Backend: Golang vs Hono.js

**Rekomendasi: Golang untuk API utama**

| Aspek | Golang | Hono.js (Node) |
|-------|--------|----------------|
| Memory footprint | ~15–20 MB | ~100–150 MB |
| Performa concurrency | Sangat tinggi (goroutines) | Baik (event loop) |
| Deploy di VPS 2GB | ✅ Sangat efisien | ⚠️ Lebih boros RAM |
| Familiar bagi Anas | ✅ Ya | ✅ Ya |
| Shared types dengan frontend | ❌ (perlu OpenAPI → TS gen) | ✅ (Hono RPC native) |
| Binary deployment | ✅ Single binary | ❌ Perlu Node runtime |

**Alasan utama memilih Golang:** VPS 2GB RAM sangat terbatas. Python CREPE analyzer sendiri membutuhkan ~300–400MB. Jika Node.js dipakai untuk API utama, sisa RAM untuk PostgreSQL dan Python service menjadi sangat tipis. Golang binary yang ringan (~20MB RAM) memberi ruang lebih besar untuk komponen lain.

*Untuk type safety dengan frontend:* Gunakan openapi-generator untuk generate TypeScript types dari Golang API schema (via huma atau swaggo).

---

### 7.2 Stack Teknologi

| Layer | Teknologi |
|-------|-----------|
| **Frontend** | React 19 · Vite 8 · TanStack Router (file-based) · TanStack Query · Tailwind CSS v4 |
| **UI Components** | Lucide React (icon library) |
| **API Client** | Orval — auto-generate React Query hooks dari Swagger spec |
| **Backend API** | Golang Echo v4 |
| **ORM** | GORM + PostgreSQL |
| **API Docs** | Swagger (swaggo/swag) — auto-generate dari annotations |
| **Audio Analyzer** | Python FastAPI + CREPE (tiny) + librosa |
| **YouTube Extraction** | yt-dlp (CLI di server) |
| **AI Explanation** | Anthropic Claude API |
| **Database** | PostgreSQL ≥ 15 |
| **Validation** | Zod (frontend) · go-playground/validator (backend) |
| **Monorepo** | pnpm workspaces · Nx (build orchestration & caching) |
| **Linting** | Biome · Husky (pre-commit) |
| **Hot Reload** | Air (backend Golang) |
| **Logging** | Logrus (backend) |
| **Containerization** | Docker + Docker Compose |

#### Workflow Type Safety (Frontend ↔ Backend)

Tidak ada manual type sharing. Tipe TypeScript di frontend di-generate otomatis dari Swagger spec backend:

```
Golang handler + swaggo annotations
        ↓
  swagger.json (auto-generated)
        ↓
  Orval codegen
        ↓
  src/api/endpoints/  ← React Query hooks (typed)
  src/api/model/      ← TypeScript types
```

Setiap perubahan schema di backend cukup jalankan:
```bash
pnpm --filter platform generate:api
```

---

### 7.3 Monorepo Structure

Mengikuti struktur dari starter kit `monorepo_gots_starterkit` dengan penambahan `apps/analyzer` untuk Python service.

```
maqam-detector/
├── apps/
│   ├── platform/               # React SPA (TanStack Router + Vite)
│   │   └── src/
│   │       ├── routes/         # File-based routing (TanStack Router)
│   │       │   ├── __root.tsx
│   │       │   ├── index.tsx   # Home / Landing
│   │       │   ├── result.$analysisId.tsx
│   │       │   └── history.tsx
│   │       ├── features/
│   │       │   ├── analyzer/   # Input panel, processing screen
│   │       │   ├── result/     # Hasil analisis, maqam card
│   │       │   └── history/    # Riwayat analisis
│   │       ├── components/
│   │       │   ├── atoms/      # Button, Input, Badge, dll
│   │       │   └── molecules/  # MaqamResultCard, HistoryItem, dll
│   │       └── api/
│   │           ├── endpoints/  # Auto-generated React Query hooks (Orval)
│   │           ├── model/      # Auto-generated TypeScript types (Orval)
│   │           └── mutator/
│   │               └── custom-instance.ts
│   ├── api/                    # Golang Echo API
│   │   ├── config/             # DB connection, env
│   │   ├── model/              # GORM models
│   │   ├── dto/                # Request/Response DTOs + Swagger annotations
│   │   ├── repository/         # Data access layer
│   │   ├── service/            # Business logic
│   │   ├── handler/            # HTTP handlers (Echo)
│   │   ├── middleware/         # CORS, logging, rate limiter
│   │   ├── docs/               # Auto-generated Swagger docs (swaggo)
│   │   └── main.go
│   └── analyzer/               # Python FastAPI (pitch extraction)
│       ├── app/
│       │   ├── routes/
│       │   ├── services/
│       │   │   ├── pitch_extractor.py   # CREPE tiny + librosa
│       │   │   ├── maqam_matcher.py     # Template matching
│       │   │   └── youtube_fetcher.py   # yt-dlp wrapper
│       │   └── data/
│       │       └── maqam_templates.json # 8 maqam core templates
│       ├── requirements.txt
│       └── Dockerfile
├── docs/                       # Dokumentasi produk (PRD, API Contract, dll)
├── nx.json                     # Nx build orchestrator
├── pnpm-workspace.yaml
├── package.json
├── docker-compose.yml
└── .env
```

> `apps/analyzer` (Python) dikelola terpisah dari ekosistem pnpm/Nx dan dideploy via Docker. `apps/api` dan `apps/platform` dikelola Nx dengan task caching.

---

### 7.4 API Endpoints (Golang)

```
POST   /api/v1/analyze/youtube   → Analisis dari URL YouTube
POST   /api/v1/analyze/upload    → Analisis dari file upload
POST   /api/v1/analyze/record    → Analisis dari rekaman mikrofon
GET    /api/v1/history           → Riwayat analisis user
GET    /api/v1/maqamat           → Daftar semua maqam yang didukung
POST   /api/v1/feedback          → Koreksi hasil deteksi
GET    /api/v1/admin/stats       → Statistik (admin only)
```

---

### 7.5 Alur Data Analisis YouTube

```
1. POST /api/v1/analyze/youtube { url: "https://youtube.com/..." }
2. Golang API: validasi URL (format + domain)
3. Golang API: POST ke Python Analyzer /internal/analyze { url }
4. Python: yt-dlp → download audio-only (~3–8MB, temp file)
5. Python: CREPE tiny → extract pitch per frame (Hz values)
6. Python: Konversi ke 24-TET Pitch Class Profile (24 bins)
7. Python: Cosine similarity vs semua maqam templates
8. Python: Return top 3 maqam + confidence scores
9. Golang API: POST ke Claude API dengan prompt + hasil maqam
10. Claude API: Generate penjelasan dalam Bahasa Indonesia
11. Golang API: Simpan ke PostgreSQL (history)
12. Golang API: Hapus temp audio file
13. Golang API: Return response ke frontend
```

---

### 7.6 Maqam Template Structure (JSON)

```json
{
  "id": "hijaz",
  "name_arabic": "حجاز",
  "name_latin": "Hijaz",
  "tonic": "D",
  "intervals_semitones": [0, 1, 5, 7, 8, 11, 12],
  "intervals_quarter_tone": [0, 2, 10, 14, 16, 22, 24],
  "characteristic_notes": ["Eb", "F#"],
  "emotion": ["dramatis", "kerinduan", "agung"],
  "common_songs": ["Ya Hanana", "Mast Qalandar variations"],
  "pitch_class_profile": [1.0, 0.8, 0.0, 0.0, 0.6, ...]
}
```

---

## 8. Non-Functional Requirements

### 8.1 Performance

| Metrik | Target |
|--------|--------|
| Waktu analisis YouTube (3–5 mnt) | ≤ 30 detik |
| Waktu analisis file upload | ≤ 20 detik |
| Waktu analisis rekaman mikrofon | ≤ 10 detik |
| Response time API (non-analisis) | ≤ 500ms |
| Concurrent users (phase 1) | ≥ 10 simultan |

### 8.2 Catatan Akurasi: Vokal & Instrumen

Sistem mendeteksi pitch dari **melodi dominan** dalam audio. Secara teknis, audio dari grup vokal + rebana (tanpa instrumen melodi) memberikan kondisi paling optimal karena rebana sebagai perkusi tidak memiliki frekuensi pitch yang dapat mengganggu analisis. Hal ini perlu dikomunikasikan kepada pengguna melalui UI agar ekspektasi terhadap hasil deteksi sesuai dengan kondisi audio input mereka.

Catatan khusus yang perlu ditampilkan di UI:
- Jika confidence < 70%: tampilkan saran "Coba dengan rekaman yang lebih jelas atau segmen melodi yang lebih panjang"
- Untuk input multi-vokal/koor: tampilkan disclaimer bahwa sistem mengambil pitch yang paling dominan

### 8.3 Reliability

- Uptime target: ≥ 99% (fase internal)
- Graceful error handling: setiap error dikembalikan dengan pesan yang human-readable
- Temp file cleanup: audio temp selalu dihapus setelah proses, bahkan jika terjadi error

### 8.3 Security

- Rate limiting: maks. 10 analisis per user per jam (fase internal)
- File validation: validasi MIME type dan magic bytes, bukan hanya ekstensi
- Temp file isolation: disimpan di direktori terisolasi dengan nama random (UUID)
- HTTPS wajib (SSL/TLS via Let's Encrypt)
- Environment secrets tidak di-commit ke repository

### 8.4 Privacy

- File audio yang di-upload atau didownload dari YouTube **tidak disimpan permanen**
- Dihapus segera setelah analisis selesai atau gagal
- Tidak ada rekaman audio yang dikirim ke pihak ketiga selain kebutuhan analisis

### 8.5 Accessibility & UX

- Mobile-friendly (responsive design)
- Bahasa antarmuka: Bahasa Indonesia
- Loading state yang jelas untuk setiap proses async
- Error messages yang informatif dan actionable

### 8.6 Infrastruktur (VPS 2 Core, 2GB RAM, 20GB Storage)

| Service | RAM Allocation |
|---------|---------------|
| OS + overhead | ~300 MB |
| PostgreSQL | ~200 MB |
| Golang API | ~20 MB |
| Python FastAPI + CREPE tiny | ~400 MB |
| TanStack SSR (Node.js) | ~200 MB |
| Buffer/cache | ~880 MB |

---

## 9. Assumptions & Constraints

### Assumptions

- Pengguna fase internal sudah familiar dengan konsep maqam secara umum meski tidak hafal teorinya
- Lagu-lagu referensi yang digunakan adalah lagu Arab/Timur Tengah dengan maqam yang dominan (bukan genre campuran)
- Koneksi internet di VPS stabil untuk mengakses YouTube dan Claude API
- Komunitas bersedia memberikan feedback dan koreksi untuk membantu meningkatkan akurasi

### Constraints

- **Legal:** yt-dlp untuk download audio dari YouTube berada di area abu-abu Terms of Service. Penggunaan dibatasi untuk keperluan analisis non-komersial dan internal komunitas
- **Akurasi:** Deteksi maqam berbasis pitch profile memiliki keterbatasan pada lagu dengan modulasi maqam yang sering atau ornamentasi yang sangat kompleks
- **Resource VPS:** CREPE harus menggunakan model `tiny` (bukan `full`) karena keterbatasan RAM
- **Maqam coverage:** Fase awal hanya mendukung 8 maqam utama; maqam yang jarang akan diproses dengan akurasi lebih rendah

---

## 10. Open Questions

| # | Pertanyaan | Status | Keputusan |
|---|-----------|--------|-----------|
| OQ-001 | Apakah sistem perlu mendukung maqam Turki (Makam) dan Persia (Dastgah) di masa depan? | ✅ Resolved | **Bisa didukung**, namun sangat tidak prioritas. Masuk roadmap jangka panjang setelah semua maqam Arab utama stabil. |
| OQ-002 | Apakah ada koleksi audio berlabel maqam dari komunitas yang bisa digunakan untuk validasi awal? | ✅ Resolved | **Tidak ada.** Komunitas mendapatkan referensi lagu dari pencarian YouTube atau channel langganan berdasarkan nama penyanyi/grup, lalu explore secara mandiri. Validasi akurasi akan dilakukan secara organik via feedback komunitas (F-013), bukan dari dataset berlabel. |
| OQ-003 | Siapa yang berperan sebagai "expert reviewer" internal untuk memvalidasi hasil koreksi komunitas? | ✅ Resolved | **Tidak diperlukan untuk saat ini.** Fitur koreksi komunitas (F-013) berjalan self-moderated di fase internal. Expert reviewer baru dipertimbangkan jika aplikasi sudah publik dan volume koreksi meningkat signifikan. |
| OQ-004 | Apakah perlu fitur notifikasi (WhatsApp/Telegram) untuk hasil analisis yang sudah selesai? | ✅ Resolved | **Tidak diperlukan.** Cukup dengan loading indicator real-time di halaman analisis. |
| OQ-005 | Untuk fase publik, apakah model monetisasi yang dipertimbangkan (freemium, donasi, atau sepenuhnya gratis)? | ✅ Resolved | **Donasi (donation-based).** Aplikasi tetap gratis untuk semua pengguna; biaya operasional (server, API) ditopang dari donasi sukarela komunitas. |

---

## 11. References

### Musik & Teori Maqam

- **MaqamWorld** — Referensi komprehensif sistem maqam Arab dengan notasi dan audio: https://www.maqamworld.com
- **MaqamLessons** — Analisis maqam dan struktur jins: https://maqamlessons.com
- **Inside Arabic Music** — Sami Abu Shumays & Scott Marcus (Oxford University Press, 2019)
- **Wikipedia: Arabic Maqam** — https://en.wikipedia.org/wiki/Arabic_maqam

### Teknologi & Research

- **CREPE: A Convolutional Representation for Pitch Estimation** — Kim et al. (ICASSP 2018): https://github.com/marl/crepe
- **MORTY: A Toolbox for Mode Recognition and Tonic Identification** — Karakurt et al. (2016)
- **MTG Makam Recognition Dataset** — Music Technology Group, Universitat Pompeu Fabra: https://github.com/MTG/otmm_makam_recognition_dataset
- **librosa: Audio and Music Signal Analysis in Python** — https://librosa.org
- **yt-dlp** — YouTube audio extraction: https://github.com/yt-dlp/yt-dlp

### Tech Stack

- **TanStack Start** — https://tanstack.com/start
- **Huma (Golang API framework)** — https://huma.rocks
- **Anthropic Claude API** — https://docs.anthropic.com
- **FastAPI** — https://fastapi.tiangolo.com

---

*Dokumen ini bersifat living document dan akan diperbarui seiring perkembangan proyek.*

**Versi History:**

| Versi | Tanggal | Perubahan |
|-------|---------|-----------|
| 1.0.3 | Mei 2026 | Revisi section 7.2 stack teknologi & 7.3 monorepo structure (Echo, GORM, Orval, Nx, Lucide React) |
| 1.0.2 | Mei 2026 | Tambah section 3.1 Deteksi Vokal vs Instrumen + NFR 8.2 |
| 1.0.1 | Mei 2026 | Resolusi semua Open Questions (OQ-001 s/d OQ-005) |
| 1.0.0 | Mei 2026 | Draft awal |

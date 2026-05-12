# UX Flow & UI Specification — MaqamDetector
**Versi:** 1.1.0  
**Tanggal:** Mei 2026  
**Target:** Figma MCP  
**Scope:** Core Features F-001 s/d F-007

---

## Daftar Isi

1. [Design Tokens](#1-design-tokens)
2. [Layout & Grid](#2-layout--grid)
3. [UX Flow Overview](#3-ux-flow-overview)
4. [Screen: Landing / Home](#4-screen-landing--home)
5. [Screen: Processing (Loading)](#5-screen-processing-loading)
6. [Screen: Hasil Analisis (F-005 & F-006)](#6-screen-hasil-analisis-f-005--f-006)
7. [Screen: Riwayat Analisis (F-007)](#7-screen-riwayat-analisis-f-007)
8. [Screen: Error States](#8-screen-error-states)
9. [Komponen Bersama (Shared Components)](#9-komponen-bersama-shared-components)

---

## 1. Design Tokens

### 1.1 Warna

Semua token tersedia dalam dua mode: **Light** dan **Dark**.

#### Brand / Primary
Terinspirasi dari estetika kaligrafi Arab — nada teal/hijau yang merepresentasikan kedalaman spiritual.

| Token | Light (Hex) | Dark (Hex) | Kegunaan |
|-------|------------|-----------|---------|
| `color.brand.primary` | `#0F6E56` | `#1D9E75` | CTA utama, link aktif |
| `color.brand.primary.hover` | `#085041` | `#5DCAA5` | Hover state CTA |
| `color.brand.primary.subtle` | `#E1F5EE` | `#04342C` | Background badge, highlight |
| `color.brand.secondary` | `#534AB7` | `#7F77DD` | Aksen sekunder, emotion tag |
| `color.brand.secondary.subtle` | `#EEEDFE` | `#26215C` | Background emotion tag |

#### Neutral
| Token | Light (Hex) | Dark (Hex) | Kegunaan |
|-------|------------|-----------|---------|
| `color.neutral.900` | `#2C2C2A` | `#F1EFE8` | Teks utama |
| `color.neutral.700` | `#444441` | `#D3D1C7` | Teks sekunder |
| `color.neutral.500` | `#5F5E5A` | `#B4B2A9` | Placeholder, hint |
| `color.neutral.300` | `#B4B2A9` | `#5F5E5A` | Border, divider |
| `color.neutral.100` | `#D3D1C7` | `#444441` | Border subtle |
| `color.neutral.50` | `#F1EFE8` | `#2C2C2A` | Background surface |
| `color.neutral.0` | `#FFFFFF` | `#1a1a18` | Background card |

#### Semantic
| Token | Light (Hex) | Dark (Hex) | Kegunaan |
|-------|------------|-----------|---------|
| `color.success` | `#3B6D11` | `#97C459` | Confidence tinggi, sukses |
| `color.success.subtle` | `#EAF3DE` | `#173404` | Background success |
| `color.warning` | `#854F0B` | `#EF9F27` | Confidence sedang, peringatan |
| `color.warning.subtle` | `#FAEEDA` | `#412402` | Background warning |
| `color.error` | `#A32D2D` | `#E24B4A` | Error, confidence rendah |
| `color.error.subtle` | `#FCEBEB` | `#501313` | Background error |
| `color.info` | `#185FA5` | `#378ADD` | Info, sedang proses |
| `color.info.subtle` | `#E6F1FB` | `#042C53` | Background info |

#### Confidence Score Colors
| Label | Range | Token | Light | Dark |
|-------|-------|-------|-------|------|
| Sangat Tinggi | ≥ 90% | `color.confidence.very_high` | `#3B6D11` | `#97C459` |
| Tinggi | 75–89% | `color.confidence.high` | `#639922` | `#C0DD97` |
| Sedang | 60–74% | `color.confidence.medium` | `#854F0B` | `#EF9F27` |
| Rendah | 40–59% | `color.confidence.low` | `#993C1D` | `#F0997B` |
| Sangat Rendah | < 40% | `color.confidence.very_low` | `#A32D2D` | `#E24B4A` |

---

### 1.2 Tipografi

Font utama: **Inter** (fallback: system-ui, sans-serif)  
Font Arab: **Amiri** (untuk nama maqam dalam huruf Arab)

| Token | Size | Weight | Line Height | Kegunaan |
|-------|------|--------|-------------|---------|
| `type.display.md` | 36px | 700 | 1.2 | Heading hero |
| `type.heading.xl` | 28px | 600 | 1.3 | Page title, nama maqam di result |
| `type.heading.lg` | 22px | 600 | 1.3 | Card title, section heading |
| `type.heading.md` | 18px | 600 | 1.4 | Sub-section title |
| `type.heading.sm` | 16px | 600 | 1.4 | Label besar |
| `type.body.lg` | 16px | 400 | 1.6 | Body teks utama |
| `type.body.md` | 14px | 400 | 1.6 | Body teks sekunder |
| `type.body.sm` | 12px | 400 | 1.5 | Caption, helper text |
| `type.label.lg` | 14px | 500 | 1.4 | Button, form label |
| `type.label.md` | 12px | 500 | 1.4 | Badge, chip, tag |
| `type.mono` | 13px | 400 | 1.5 | Interval tangga nada, URL |
| `type.arabic` | 24px | 400 | 1.6 | Nama maqam Arab — font: Amiri |

---

### 1.3 Spacing

Base grid: **8px**

| Token | Value | Kegunaan |
|-------|-------|---------|
| `space.1` | 4px | Gap micro |
| `space.2` | 8px | Gap kecil, inner padding kompak |
| `space.3` | 12px | Padding internal komponen kecil |
| `space.4` | 16px | Padding standar, gap antar elemen |
| `space.5` | 20px | Padding komponen medium |
| `space.6` | 24px | Padding card, section gap |
| `space.8` | 32px | Gap antar section |
| `space.10` | 40px | Padding section besar |
| `space.12` | 48px | Margin vertikal section |
| `space.16` | 64px | Hero padding |

---

### 1.4 Border Radius

| Token | Value | Kegunaan |
|-------|-------|---------|
| `radius.sm` | 4px | Badge, chip |
| `radius.md` | 8px | Input, button |
| `radius.lg` | 12px | Card, panel |
| `radius.xl` | 16px | Card besar, sheet |
| `radius.2xl` | 24px | Bottom sheet, modal |
| `radius.full` | 9999px | Pill button, progress bar |

---

### 1.5 Shadow (Elevation)

| Token | CSS Value | Kegunaan |
|-------|-----------|---------|
| `shadow.sm` | `0 1px 3px rgba(0,0,0,0.08)` | Input focus |
| `shadow.md` | `0 4px 12px rgba(0,0,0,0.10)` | Card, dropdown |
| `shadow.lg` | `0 8px 24px rgba(0,0,0,0.12)` | Modal, bottom sheet |
| `shadow.brand` | `0 4px 16px rgba(15,110,86,0.25)` | CTA button hover |

---

### 1.6 Ikon

Library: **Lucide React** (default size 24px, stroke-width 1.5)

| Konteks | Nama Lucide |
|---------|-------------|
| YouTube input | `Youtube` |
| Upload file | `UploadCloud` |
| Mikrofon / rekam | `Mic` |
| Humming | `AudioWaveform` |
| Stop rekam | `StopCircle` |
| Rekam ulang | `RotateCcw` |
| Sukses / check | `CheckCircle2` |
| Error | `XCircle` |
| Peringatan | `AlertTriangle` |
| Info | `Info` |
| Riwayat | `History` |
| Hapus | `Trash2` |
| Musik / maqam | `Music` |
| Close | `X` |
| Kembali | `ArrowLeft` |
| Loading | `Loader2` (animated spin) |
| Context menu | `MoreVertical` |
| Musik note (list) | `Music2` |

---

## 2. Layout & Grid

### 2.1 Breakpoint

| Nama | Min Width |
|------|-----------|
| `mobile` | 0px |
| `tablet` | 768px |
| `desktop` | 1024px |
| `wide` | 1280px |

### 2.2 Container

| Breakpoint | Padding horizontal | Max width konten |
|------------|-------------------|-----------------|
| Mobile | 16px | 100% |
| Tablet | 24px | 720px |
| Desktop | 32px | 960px |
| Wide | 40px | 1120px |

### 2.3 Frame Sizes (Figma)

| Frame | Width | Height |
|-------|-------|--------|
| `Mobile/Home` | 390px | Auto |
| `Mobile/Processing` | 390px | Auto |
| `Mobile/Result` | 390px | Auto |
| `Mobile/History` | 390px | Auto |
| `Desktop/Home` | 1440px | Auto |
| `Desktop/Processing` | 1440px | Auto |
| `Desktop/Result` | 1440px | Auto |
| `Desktop/History` | 1440px | Auto |

---

## 3. UX Flow Overview

### 3.1 Master Flow

```
[Home / Landing]
       │
       ├── Tab YouTube  (F-001)
       ├── Tab Upload   (F-002)
       └── Tab Rekam    (F-003 & F-004)
                │
                ▼
      [Processing / Loading]
                │
        ┌───────┴────────┐
        ▼                ▼
  [Hasil Analisis]   [Error State]
  (F-005 & F-006)        │
        │            [Coba Lagi / Home]
        ▼
  [Analisis Lagi] → kembali ke Home

  [Riwayat] (F-007) ← dari navbar, kapan saja
```

### 3.2 Happy Path (YouTube)

| Step | Screen | Aksi User |
|------|--------|-----------|
| 1 | Home | Pilih tab YouTube |
| 2 | Home | Paste URL YouTube |
| 3 | Home | Tap "Analisis Maqam" |
| 4 | Processing | Menunggu (10–30 detik) |
| 5 | Hasil | Baca nama maqam + penjelasan |
| 6 | Hasil | Tap "Analisis Lagu Lain" → kembali Home |

---

## 4. Screen: Landing / Home

### 4.1 Struktur Halaman

```
┌─────────────────────────────────┐
│  [Navbar]                       │
├─────────────────────────────────┤
│  [Hero Section]                 │
│  Judul + Deskripsi singkat      │
├─────────────────────────────────┤
│  [Input Method Tabs]            │
│  YouTube | Upload | Rekam       │
├─────────────────────────────────┤
│  [Input Panel]                  │
│  (konten berubah sesuai tab)    │
├─────────────────────────────────┤
│  [CTA: Analisis Maqam]          │
├─────────────────────────────────┤
│  [Maqam yang Didukung]          │
│  8 maqam dalam chip row         │
└─────────────────────────────────┘
```

---

### 4.2 Navbar

**Height:** 56px · **Sticky:** Ya  
**Background:** `color.neutral.0` · **Border bottom:** `1px solid color.neutral.100`

| | Kiri | Kanan |
|--|------|-------|
| Desktop | Logo + wordmark "MaqamDetector" | Link "Riwayat" |
| Mobile | Logo + wordmark | Ikon `History` |

**Logo:** 32px × 32px · warna `color.brand.primary` · wordmark `type.heading.sm`

---

### 4.3 Hero Section

**Padding:** `space.12` atas · `space.8` bawah · **Alignment:** Center

| Elemen | Spec |
|--------|------|
| Judul | `type.display.md` · `color.neutral.900` · "Kenali Maqam Lagu Arab & Timur Tengah" |
| Sub-judul | `type.body.lg` · `color.neutral.700` · max-width 540px · "Upload lagu, tempel link YouTube, atau nyanyikan melodinya — kami bantu identifikasi maqam-nya dalam detik." |
| Gap | `space.4` |

---

### 4.4 Input Method Tabs

**Komponen:** `InputMethodTabs`  
**Style:** Segmented control / tab underline

```
[🎬 YouTube]   [📁 Upload File]   [🎤 Rekam / Humming]
```

| State | Border bottom | Text color |
|-------|--------------|-----------|
| Active | `2px solid color.brand.primary` | `color.brand.primary` |
| Inactive | — | `color.neutral.500` |
| Hover | `2px solid color.neutral.300` | `color.neutral.700` |

- **Mobile:** 3 tab rata lebar, ikon di atas label
- **Desktop:** auto width, ikon inline dengan label

---

### 4.5 Input Panel — Tab YouTube (F-001)

**Komponen:** `YoutubeInputPanel`  
**Background:** `color.neutral.50` · **Border:** `1px solid color.neutral.100` · **Border radius:** `radius.xl` · **Padding:** `space.6`

```
┌─────────────────────────────────────────┐
│  Label: "Link YouTube"                  │
│  ┌─────────────────────────────────────┐│
│  │ 🔗 Paste link YouTube di sini...   ││
│  └─────────────────────────────────────┘│
│                                         │
│  [Video Preview Card — muncul jika URL valid]
│                                         │
│  Label: "Segmen yang dianalisis"        │
│  [Mulai: 0 detik]  [Durasi: 60 detik]  │
└─────────────────────────────────────────┘
```

#### URL Input Field

| Prop | Nilai |
|------|-------|
| Height | 48px |
| Border radius | `radius.md` |
| Placeholder | "https://www.youtube.com/watch?v=..." |
| Icon kiri | `Youtube` 20px `color.neutral.500` |

| State | Border | Background | Keterangan |
|-------|--------|-----------|-----------|
| Default | `color.neutral.300` | `color.neutral.0` | |
| Focus | `color.brand.primary` + `shadow.sm` | `color.neutral.0` | |
| Valid | `color.success` | `color.success.subtle` | Icon ✓ di kanan |
| Error | `color.error` | `color.error.subtle` | Helper text muncul |

#### Video Preview Card

Muncul setelah URL valid (debounce 800ms).

| Elemen | Spec |
|--------|------|
| Container | `radius.md` · border `1px solid color.neutral.100` · padding `space.3` |
| Thumbnail | 80px × 56px · `radius.sm` · object-fit cover |
| Judul video | `type.body.md` / 600 · 2 baris max · truncate |
| Meta | `type.body.sm` · `color.neutral.500` (durasi · channel) |
| Loading | Skeleton shimmer |

#### Segment Selector

| Elemen | Spec |
|--------|------|
| Input "Mulai" | 80px · type number · suffix "detik" · default 0 |
| Input "Durasi" | 80px · type number · suffix "detik" · default 60 · max 120 |
| Helper | `type.body.sm` · "Sistem menganalisis 60 detik pertama secara default" |

---

### 4.6 Input Panel — Tab Upload File (F-002)

**Komponen:** `FileUploadPanel`

```
┌─────────────────────────────────────────┐
│                                         │
│   [ ↑ ]                                 │
│   Seret file ke sini atau               │
│   [Pilih File dari Perangkat]           │
│                                         │
│   MP3, WAV, M4A, FLAC, OGG · Maks 50MB │
│                                         │
└─────────────────────────────────────────┘
```

**Min-height:** 160px · **Border:** `2px dashed color.neutral.300` · **Border radius:** `radius.xl`

| State | Border | Background |
|-------|--------|-----------|
| Default | `2px dashed color.neutral.300` | `color.neutral.50` |
| Drag over | `2px dashed color.brand.primary` | `color.brand.primary.subtle` |
| File terpilih | `2px solid color.neutral.200` | `color.neutral.0` |
| Error format | `2px dashed color.error` | `color.error.subtle` |

**Setelah file terpilih** → dropzone berubah jadi **File Info Card:**

```
┌─────────────────────────────────────────┐
│ 🎵  ya_hanana.mp3                   [×] │
│     4.6 MB · 3:34                       │
└─────────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Ikon | `Music` 24px `color.brand.primary` |
| Nama file | `type.body.md` / 500 · truncate tengah |
| Meta | `type.body.sm` `color.neutral.500` |
| Tombol hapus | `X` ikon 16px · tap area 32px × 32px |

---

### 4.7 Input Panel — Tab Rekam / Humming (F-003 & F-004)

**Komponen:** `RecordPanel`

```
┌─────────────────────────────────────────┐
│  Mode:  ○ Mikrofon   ● Humming         │
│                                         │
│  [     Waveform Visualizer Area     ]   │
│                                         │
│            [ 🎤 Mulai Rekam ]          │
│                                         │
│  Senandungkan melodi utama lagu         │
│  minimal 5 detik                        │
└─────────────────────────────────────────┘
```

#### Mode Toggle

| State | Style |
|-------|-------|
| Active | `color.brand.primary` fill · label `type.label.lg` / 600 |
| Inactive | Border `color.neutral.300` · label `type.body.md` |

#### Waveform Visualizer

| State | Tampilan |
|-------|---------|
| Idle | Area kosong · `2px dashed color.neutral.200` |
| Recording | Bar animasi warna `color.brand.primary` (WebAudio API) |
| Done | Waveform static · durasi di pojok kanan |

**Height:** 80px · **Background:** `color.neutral.50` · **Border radius:** `radius.md`

#### Record Button

| State | Label | Background | Ikon |
|-------|-------|-----------|------|
| Idle | "Mulai Rekam" | `color.brand.primary` | `Mic` |
| Recording | "Berhenti · 0:00" | `color.error` | `StopCircle` (pulse) |
| Done | "Rekam Ulang" | `color.neutral.50` + border | `RotateCcw` |

**Sub-state Done:**
```
┌─────────────────────────────────────────┐
│  Mode:  ○ Mikrofon   ● Humming         │
│                                         │
│  ████ ██ █████ ███ ████ ██ █████   0:23 │  ← waveform static
│                                         │
│         [↺ Rekam Ulang]                │
└─────────────────────────────────────────┘

[Analisis Maqam]  ← enabled
```

---

### 4.8 CTA Button — Analisis Maqam

**Komponen:** `AnalyzeButton`  
**Width:** Full width (mobile) · auto min-width 200px (desktop)  
**Height:** 52px · **Border radius:** `radius.md` · **Font:** `type.label.lg`

| State | Background | Text | Shadow |
|-------|-----------|------|--------|
| Enabled | `color.brand.primary` | `#FFFFFF` | — |
| Hover | `color.brand.primary.hover` | `#FFFFFF` | `shadow.brand` |
| Disabled | `color.neutral.200` | `color.neutral.500` | — |
| Loading | `color.brand.primary` 50% | — | Spinner tengah |

**Kapan disabled:**
- YouTube: URL belum valid
- Upload: belum ada file
- Rekam: durasi rekaman = 0

---

### 4.9 Maqam yang Didukung

```
Maqam yang didukung:
[Hijaz] [Rast] [Bayati] [Nahawand] [Kurd] [Saba] [Ajam] [Jiharkah]
```

| Elemen | Spec |
|--------|------|
| Label | `type.body.sm` · `color.neutral.500` |
| Chip | `radius.full` · padding `space.2` `space.4` · border `1px solid color.neutral.200` |
| Chip text | `type.label.md` · `color.neutral.700` |
| Chip hover | Background `color.neutral.50` · border `color.brand.primary` |
| Layout | Horizontal scroll (mobile) · wrap (desktop) |

---

## 5. Screen: Processing (Loading)

### 5.1 Struktur

Navigasi ke screen ini segera setelah user tap "Analisis Maqam".

```
┌─────────────────────────────────────────┐
│  [Navbar]                               │
├─────────────────────────────────────────┤
│                                         │
│   [Info Sumber Audio]                   │
│   Thumbnail / ikon + judul / nama file  │
│                                         │
│   ──────────────────────────────────    │
│                                         │
│   [Progress Steps]                      │
│   ✓  Mengunduh audio                    │
│   ◌  Menganalisis pitch                 │  ← active
│   ○  Mencocokkan pola maqam             │
│   ○  Menyusun penjelasan                │
│                                         │
│   ──────────────────────────────────    │
│   Estimasi: sekitar 15 detik lagi...    │
│                                         │
│   [Batalkan]                            │
│                                         │
└─────────────────────────────────────────┘
```

---

### 5.2 Komponen: AnalysisProgressSteps

| State | Ikon | Warna ikon | Teks |
|-------|------|-----------|------|
| Pending | `Circle` outline | `color.neutral.300` | `color.neutral.400` |
| Active | `Loader2` animasi | `color.brand.primary` | `color.neutral.900` |
| Completed | `CheckCircle2` filled | `color.success` | `color.neutral.700` |

Step labels:
1. "Mengunduh audio" *(hanya tampil untuk input YouTube)*
2. "Menganalisis pitch melodi"
3. "Mencocokkan pola maqam"
4. "Menyusun penjelasan"

**Layout:** Vertical list · gap `space.4` · connector line `1px solid color.neutral.200`

---

### 5.3 Info Sumber Audio

| Input type | Tampilan |
|------------|---------|
| YouTube | Thumbnail 80px × 56px + judul + channel |
| Upload | Ikon `Music` + nama file + durasi |
| Mikrofon | Ikon `Mic` + "Rekaman mikrofon · 0:23" |
| Humming | Ikon `AudioWaveform` + "Senandung · 0:18" |

---

## 6. Screen: Hasil Analisis (F-005 & F-006)

### 6.1 Struktur Desktop (2 kolom)

```
┌─────────────────────────────────────────┐
│  [Navbar]                               │
├──────────────┬──────────────────────────┤
│              │                          │
│  [Info       │  [Maqam Result Card]     │
│   Sumber]    │                          │
│              │  [Kandidat Alternatif]   │
│  [Kandidat   │                          │
│   Lain]      │  [Penjelasan Detail]     │
│              │  - Karakteristik & Emosi │
│              │  - Struktur Tangga Nada  │
│              │  - Contoh Lagu           │
│              │  - Tips Banjari          │
│              │                          │
│              │  [Analisis Lagu Lain]    │
└──────────────┴──────────────────────────┘
```

### 6.2 Struktur Mobile (1 kolom, scroll)

```
[Navbar + Tombol Kembali]
[Info Sumber Audio]
[Maqam Result Card]
[Confidence & Kandidat Alternatif]
[Penjelasan Detail]
  - Karakteristik & Emosi
  - Struktur Tangga Nada
  - Contoh Lagu
  - Tips Aransemen Banjari
[Analisis Lagu Lain]
```

---

### 6.3 Komponen: MaqamResultCard

**Background:** Gradient `color.brand.primary.subtle` → `color.neutral.0`  
**Border:** `1px solid color.brand.primary` (20% opacity)  
**Border radius:** `radius.xl` · **Padding:** `space.6`

```
┌─────────────────────────────────────────┐
│  MAQAM TERDETEKSI                       │
│                                         │
│  Hijaz                       حجاز      │
│  ─────────────────────────────────────  │
│  Confidence: ████████░░  87%  Tinggi   │
│                                         │
│  D – E♭ – F# – G – A – B♭ – C         │
└─────────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Label "MAQAM TERDETEKSI" | `type.label.md` · `color.brand.primary` · letter-spacing 0.08em |
| Nama Latin | `type.heading.xl` · `color.neutral.900` |
| Nama Arab | `type.arabic` 24px · `color.neutral.700` · text-align right · font Amiri |
| Divider | `1px solid color.neutral.200` |
| Progress bar | Height 8px · `radius.full` · fill sesuai confidence color token |
| Persentase | `type.heading.sm` · warna sesuai confidence token |
| Label verbal | Badge `radius.sm` · warna sesuai confidence token |
| Interval tangga nada | `type.mono` · `color.neutral.600` |

#### Confidence → Color Mapping

| Confidence | Label | Color token |
|-----------|-------|------------|
| ≥ 90% | "Sangat Tinggi" | `color.confidence.very_high` |
| 75–89% | "Tinggi" | `color.confidence.high` |
| 60–74% | "Sedang" | `color.confidence.medium` |
| 40–59% | "Rendah" | `color.confidence.low` |
| < 40% | "Sangat Rendah" | `color.confidence.very_low` |

---

### 6.4 Komponen: MaqamCandidateList

Selalu tampil (rank 2 dan 3). Emphasis lebih kuat jika confidence utama < 70%.

```
Kandidat lainnya:

┌──────────────────────────────────────┐
│  #2  Kurd               9%           │
├──────────────────────────────────────┤
│  #3  Bayati             4%           │
└──────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Label section | `type.body.sm` · `color.neutral.500` |
| Item padding | `space.3` vertical · `space.4` horizontal |
| Rank badge | `type.label.md` · `color.neutral.500` |
| Nama maqam | `type.body.md` / 500 · `color.neutral.700` |
| Persentase | `type.label.md` · warna confidence |
| Border antar item | `1px solid color.neutral.100` |

**Warning block** (jika confidence utama < 70%):

```
┌────────────────────────────────────────┐
│  ⚠️  Confidence Sedang                 │
│  Hasil ini mungkin kurang pasti.       │
│  Pertimbangkan kandidat alternatif     │
│  atau coba dengan rekaman lebih jelas. │
└────────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Background | `color.warning.subtle` |
| Border | `1px solid color.warning` 30% opacity · `radius.md` |
| Ikon | `AlertTriangle` 20px · `color.warning` |
| Teks | `type.body.sm` · `color.neutral.700` |

---

### 6.5 Komponen: MaqamExplanation

Terbagi dalam sub-section. Accordion pada mobile, tampil semua pada desktop.

#### Sub-section: Karakteristik & Emosi

```
🎭 Karakteristik & Emosi

[Teks penjelasan karakteristik maqam...]

[dramatis]  [kerinduan]  [agung]  [spiritual]
```

| Elemen | Spec |
|--------|------|
| Judul sub | `type.heading.sm` · `color.neutral.900` · ikon 20px `color.brand.primary` |
| Body teks | `type.body.md` · `color.neutral.700` · line-height 1.7 |
| Emotion tag | Chip `radius.full` · background `color.brand.secondary.subtle` · text `color.brand.secondary` · `type.label.md` |

#### Sub-section: Struktur Tangga Nada

```
🎼 Struktur Tangga Nada

D – E♭ – F# – G – A – B♭ – C

[Teks penjelasan interval khas...]
```

| Elemen | Spec |
|--------|------|
| Tangga nada | `type.mono` · background `color.neutral.50` · padding `space.3` `space.4` · `radius.md` |
| Penjelasan | `type.body.md` · `color.neutral.700` |

#### Sub-section: Contoh Lagu

```
🎵 Contoh Lagu & Munsyid

• Ya Hanana
• Qasidah Burda (pembuka)
• Tala'al Badru Alayna (versi Hijaz)
```

| Elemen | Spec |
|--------|------|
| Item | `type.body.md` · `color.neutral.700` |
| Bullet | `color.brand.primary` |
| Gap antar item | `space.2` |

#### Sub-section: Tips Aransemen Banjari

```
┌────────────────────────────────────────────┐
│  🥁  Tips untuk Aransemen Banjari          │
│                                            │
│  Maqam Hijaz sangat cocok untuk bagian     │
│  pembuka yang dramatis. Pada vokal banjari,│
│  perhatikan ornamentasi pada nada E♭...    │
└────────────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Background | `color.brand.primary.subtle` |
| Border kiri | `3px solid color.brand.primary` |
| Border radius | `radius.md` |
| Padding | `space.4` |
| Teks | `type.body.md` · `color.neutral.800` |

---

### 6.6 Action Button

Satu tombol di bawah halaman hasil:

```
[🔍 Analisis Lagu Lain]
```

| Elemen | Spec |
|--------|------|
| Style | Filled · `color.brand.primary` |
| Width | Full width (mobile) · auto (desktop) |
| Height | 52px · `radius.md` |
| Font | `type.label.lg` |
| Aksi | Navigate ke Home, reset semua input |

---

## 7. Screen: Riwayat Analisis (F-007)

### 7.1 Struktur

```
[Navbar]

[Page Title: "Riwayat Analisis"]

[Filter Bar: Semua | YouTube | Upload | Rekaman]

[List HistoryItem]

[Empty State jika kosong]
```

---

### 7.2 Komponen: HistoryItem

```
┌─────────────────────────────────────────────┐
│  [Thumb]  Judul / Sumber               [⋮]  │
│           🎵 Hijaz  ·  87%                  │
│           📅 10 Mei 2026, 14.23             │
└─────────────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Container | `radius.lg` · border `1px solid color.neutral.100` · padding `space.4` |
| Hover | Background `color.neutral.50` |
| Thumbnail / Ikon | 56px × 40px (YouTube) · 40px × 40px ikon (lainnya) · `radius.sm` |
| Judul sumber | `type.body.md` / 500 · truncate 1 baris |
| Nama maqam | `type.body.sm` · `color.neutral.700` · ikon `Music` 14px |
| Confidence | `type.body.sm` · warna confidence |
| Tanggal | `type.body.sm` · `color.neutral.500` |
| Menu [⋮] | `MoreVertical` 20px · tap area 32px × 32px |
| Status badge | Hanya tampil jika bukan `completed` — "Gagal" (error) / "Diproses" (info) |

**Context menu dari [⋮]:**
- "Lihat Detail" → navigasi ke Result screen
- "Hapus" → konfirmasi inline → `DELETE /history/:id`

---

### 7.3 Komponen: Filter Bar

```
[Semua]  [YouTube]  [Upload]  [Rekaman]
```

| State | Style |
|-------|-------|
| Active | Background `color.brand.primary.subtle` · text `color.brand.primary` · border `color.brand.primary` |
| Inactive | Border `color.neutral.200` · text `color.neutral.600` |

---

### 7.4 Empty State

```
     [ClockCounterClockwise, 64px, color.neutral.300]

          Belum ada riwayat analisis

   Mulai dengan menganalisis lagu Arab favoritmu.

            [Mulai Analisis]
```

| Elemen | Spec |
|--------|------|
| Ikon | 64px · `color.neutral.300` |
| Judul | `type.heading.md` · `color.neutral.700` |
| Sub | `type.body.md` · `color.neutral.500` · max-width 280px · center |
| CTA | Filled · `color.brand.primary` → navigate ke Home |

---

## 8. Screen: Error States

### 8.1 Error di Processing Screen

Tampil jika polling mengembalikan `status: "failed"`.

```
┌─────────────────────────────────────────┐
│                                         │
│   [XCircle, 48px, color.error]          │
│                                         │
│   Analisis Gagal                        │
│                                         │
│   [Pesan error sesuai kode]             │
│                                         │
│   [Coba Lagi]    [Kembali ke Home]      │
│                                         │
└─────────────────────────────────────────┘
```

| Elemen | Spec |
|--------|------|
| Ikon | `XCircle` 48px · `color.error` |
| Judul | `type.heading.lg` · `color.neutral.900` |
| Pesan | `type.body.md` · `color.neutral.700` · max-width 320px |
| Tombol Coba Lagi | Filled · `color.brand.primary` |
| Tombol Kembali | Ghost / text · `color.neutral.600` |

---

### 8.2 Error Mapping

| API Error Code | Pesan UI |
|----------------|---------|
| `VIDEO_UNAVAILABLE` | "Video tidak dapat diakses. Pastikan video bersifat publik dan coba lagi." |
| `VIDEO_TOO_LONG` | "Durasi video terlalu panjang. Gunakan segmen maksimal 15 menit." |
| `LOW_AUDIO_QUALITY` | "Kualitas audio terlalu rendah. Coba dengan rekaman yang lebih jelas." |
| `AUDIO_TOO_SHORT` | "Rekaman terlalu pendek. Nyanyikan melodi minimal 5 detik." |
| `ANALYSIS_FAILED` | "Terjadi kesalahan saat menganalisis. Silakan coba beberapa saat lagi." |
| `RATE_LIMIT_EXCEEDED` | "Batas analisis tercapai. Coba lagi dalam {X} menit." |
| Default | "Terjadi kesalahan. Silakan coba lagi." |

---

### 8.3 Toast Notification

Digunakan untuk feedback cepat: hapus riwayat berhasil, dll.

| Tipe | Background | Ikon | Durasi |
|------|-----------|------|--------|
| Success | `color.success.subtle` + border `color.success` | `CheckCircle2` | 3 detik |
| Error | `color.error.subtle` + border `color.error` | `XCircle` | 5 detik |
| Info | `color.info.subtle` + border `color.info` | `Info` | 3 detik |
| Warning | `color.warning.subtle` + border `color.warning` | `AlertTriangle` | 4 detik |

**Posisi:** Bottom center (mobile) · Bottom right (desktop)  
**Border radius:** `radius.md` · **Padding:** `space.3` `space.4`  
**Animation:** Slide up + fade in · auto-dismiss dengan fade out

---

## 9. Komponen Bersama (Shared Components)

### 9.1 Component Inventory (Figma Library)

| Nama Komponen | Variants |
|---------------|----------|
| `Button` | `size: sm/md/lg` · `variant: filled/outline/ghost/text` · `state: default/hover/active/disabled/loading` |
| `InputField` | `state: default/focus/valid/error/disabled` · `size: md/lg` |
| `Dropdown/Select` | `state: default/open/selected/disabled` |
| `Badge` | `variant: brand/success/warning/error/info/neutral` |
| `Chip` | `variant: outline/filled` · `state: default/hover/active` |
| `Card` | `variant: default/elevated/brand` |
| `ProgressBar` | `variant: brand/success/warning/error/medium` |
| `Spinner` | `size: sm/md/lg` · `color: brand/white/neutral` |
| `SkeletonLoader` | `variant: text/image/card` |
| `Toast` | `variant: success/error/info/warning` |
| `BottomSheet` | `state: default/expanded` *(mobile modal)* |
| `Navbar` | `variant: default` |
| `ContextMenu` | — |
| `EmptyState` | — |
| `InputMethodTabs` | `activeTab: youtube/upload/record` |
| `YoutubeInputPanel` | `state: empty/url-valid/url-error` |
| `VideoPreviewCard` | `state: loading/loaded/error` |
| `FileUploadZone` | `state: default/drag-over/selected/error` |
| `FileInfoCard` | — |
| `RecordPanel` | `mode: microphone/humming` · `state: idle/recording/done` |
| `WaveformVisualizer` | `state: idle/recording/done` |
| `RecordButton` | `state: idle/recording/done` |
| `AnalyzeButton` | `state: enabled/disabled/loading` |
| `AnalysisProgressSteps` | `activeStep: 1/2/3/4` |
| `MaqamResultCard` | `confidence: very_high/high/medium/low/very_low` |
| `MaqamCandidateList` | `showWarning: boolean` |
| `ConfidenceBadge` | `level: very_high/high/medium/low/very_low` |
| `MaqamExplanation` | — |
| `HistoryItem` | `inputType: youtube/upload/microphone/humming` · `status: completed/failed/processing` |
| `FilterBar` | `activeFilter: all/youtube/upload/record` |

---

### 9.2 Animation & Transition Specs

| Animasi | Duration | Easing | Properti |
|---------|----------|--------|---------|
| Tab switch | 200ms | ease-out | opacity, translateX |
| Card masuk | 300ms | ease-out | opacity, translateY(8px→0) |
| Bottom sheet | 250ms | ease-out | translateY |
| Toast in / out | 200ms / 150ms | ease-out / ease-in | opacity, translateY |
| Spinner | 700ms | linear infinite | rotate |
| Waveform bar | 150ms | ease-in-out | height |
| Confidence bar fill | 600ms delay 200ms | ease-out | width |
| Record pulse | 1s | ease-in-out infinite | opacity 1→0.5→1 |
| Skeleton shimmer | 1.5s | ease-in-out infinite | background-position |

---

### 9.3 Responsive Behavior Summary

| Komponen | Mobile | Desktop |
|----------|--------|---------|
| InputMethodTabs | Ikon di atas label, full width | Ikon inline label, auto width |
| Home layout | 1 kolom | 1 kolom, max-width 680px centered |
| Result layout | 1 kolom scroll | 2 kolom (sidebar + konten) |
| Modal / konfirmasi | Bottom sheet | Dialog modal centered |
| MaqamResultCard | Full width | Full width dalam kolom kanan |
| HistoryItem | Full width | Full width, max-width 760px |

---

**Versi History:**

| Versi | Tanggal | Perubahan |
|-------|---------|-----------|
| 1.1.1 | Mei 2026 | Update icon library: Phosphor → Lucide React (sesuai starter kit) |
| 1.1.0 | Mei 2026 | Hapus semua elemen F-013 (Koreksi Komunitas): modal feedback, tombol "Kirim Koreksi", komponen terkait. Fokus murni core F-001–F-007 |
| 1.0.0 | Mei 2026 | Draft awal |

*Versi 1.1.0 — Mei 2026*

# API Contract — MaqamDetector
**Base URL:** `https://api.maqamdetector.id/api/v1`  
**Versi:** 1.1.0  
**Format:** JSON (kecuali endpoint upload multipart)  
**Autentikasi:** Tidak diperlukan (fase internal — guest session via header)

---

## Konvensi Umum

### Request Headers

```
Content-Type: application/json
Accept: application/json
X-Session-ID: <uuid>   # Wajib di semua request — dibuat & disimpan di sisi client
```

> `X-Session-ID` adalah UUID yang dibuat sisi client (disimpan di localStorage) untuk mengidentifikasi sesi guest. Semua data analisis dan riwayat terikat ke session ini.

---

### Format Response Standar

**Success:**
```json
{
  "success": true,
  "data": { ... },
  "meta": { ... }
}
```

**Error:**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Pesan error yang human-readable",
    "details": [ ... ]
  }
}
```

---

### HTTP Status Codes

| Kode | Kondisi |
|------|---------|
| `200 OK` | Request berhasil |
| `202 Accepted` | Request diterima, proses async berjalan |
| `400 Bad Request` | Input tidak valid |
| `403 Forbidden` | Tidak punya akses ke resource ini |
| `404 Not Found` | Resource tidak ditemukan |
| `413 Content Too Large` | File melebihi batas ukuran |
| `415 Unsupported Media Type` | Format file tidak didukung |
| `422 Unprocessable Entity` | Request valid tapi tidak bisa diproses |
| `429 Too Many Requests` | Rate limit tercapai |
| `500 Internal Server Error` | Kesalahan server |

---

### Error Codes

| Code | Deskripsi |
|------|-----------|
| `VALIDATION_ERROR` | Field tidak valid |
| `INVALID_URL` | URL YouTube tidak valid |
| `VIDEO_UNAVAILABLE` | Video private, dihapus, atau geo-restricted |
| `VIDEO_TOO_LONG` | Durasi video melebihi batas |
| `FILE_TOO_LARGE` | Ukuran file melebihi 50MB |
| `UNSUPPORTED_FORMAT` | Format audio tidak didukung |
| `AUDIO_TOO_SHORT` | Audio terlalu pendek untuk dianalisis (< 5 detik) |
| `LOW_AUDIO_QUALITY` | Kualitas audio terlalu rendah untuk deteksi |
| `ANALYSIS_FAILED` | Proses analisis gagal di server |
| `RATE_LIMIT_EXCEEDED` | Melebihi batas 10 analisis per jam |
| `NOT_FOUND` | Resource tidak ditemukan |
| `FORBIDDEN` | Tidak punya akses ke resource ini |

---

## Endpoints

---

### 1. Analisis

#### `POST /analyze/youtube`
Meminta analisis maqam dari URL YouTube.

**Request:**
```http
POST /api/v1/analyze/youtube
Content-Type: application/json
X-Session-ID: 550e8400-e29b-41d4-a716-446655440000
```

```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "segment_start": 0,
  "segment_duration": 60
}
```

| Field | Tipe | Wajib | Deskripsi |
|-------|------|-------|-----------|
| `url` | `string` | ✅ | URL YouTube yang valid |
| `segment_start` | `integer` | ❌ | Detik mulai analisis (default: `0`) |
| `segment_duration` | `integer` | ❌ | Durasi segmen dalam detik (default: `60`, max: `120`) |

**Response `202 Accepted`:**
```json
{
  "success": true,
  "data": {
    "analysis_id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
    "status": "pending",
    "estimated_seconds": 25
  }
}
```

**Error `400`** — URL tidak valid:
```json
{
  "success": false,
  "error": {
    "code": "INVALID_URL",
    "message": "URL yang diberikan bukan URL YouTube yang valid"
  }
}
```

**Error `422`** — Video tidak bisa diakses:
```json
{
  "success": false,
  "error": {
    "code": "VIDEO_UNAVAILABLE",
    "message": "Video tidak dapat diakses. Kemungkinan video bersifat privat atau telah dihapus"
  }
}
```

**Error `429`** — Rate limit:
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Batas analisis tercapai. Coba lagi dalam 45 menit",
    "details": {
      "retry_after_seconds": 2700
    }
  }
}
```

---

#### `POST /analyze/upload`
Meminta analisis maqam dari file audio yang di-upload.

**Request:**
```http
POST /api/v1/analyze/upload
Content-Type: multipart/form-data
X-Session-ID: 550e8400-e29b-41d4-a716-446655440000
```

| Field | Tipe | Wajib | Deskripsi |
|-------|------|-------|-----------|
| `file` | `file` | ✅ | File audio (MP3, WAV, M4A, FLAC, OGG). Maks 50MB |
| `segment_start` | `integer` | ❌ | Detik mulai analisis (default: `0`) |
| `segment_duration` | `integer` | ❌ | Durasi segmen (default: `60`, max: `120`) |

**Response `202 Accepted`:**
```json
{
  "success": true,
  "data": {
    "analysis_id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
    "status": "pending",
    "estimated_seconds": 15,
    "file": {
      "name": "ya_hanana.mp3",
      "size_bytes": 4821234,
      "duration_seconds": 214
    }
  }
}
```

**Error `413`:**
```json
{
  "success": false,
  "error": {
    "code": "FILE_TOO_LARGE",
    "message": "Ukuran file melebihi batas 50MB"
  }
}
```

**Error `415`:**
```json
{
  "success": false,
  "error": {
    "code": "UNSUPPORTED_FORMAT",
    "message": "Format file tidak didukung. Gunakan MP3, WAV, M4A, FLAC, atau OGG",
    "details": {
      "received_mime": "video/mp4",
      "supported": ["audio/mpeg", "audio/wav", "audio/mp4", "audio/flac", "audio/ogg"]
    }
  }
}
```

---

#### `POST /analyze/record`
Meminta analisis maqam dari audio hasil rekaman browser (mikrofon atau humming).

**Request:**
```http
POST /api/v1/analyze/record
Content-Type: multipart/form-data
X-Session-ID: 550e8400-e29b-41d4-a716-446655440000
```

| Field | Tipe | Wajib | Deskripsi |
|-------|------|-------|-----------|
| `file` | `file` | ✅ | File audio dari browser (WebM/Opus atau WAV). Maks 10MB |
| `mode` | `string` | ✅ | `"microphone"` atau `"humming"` |
| `duration_seconds` | `integer` | ✅ | Durasi rekaman aktual dalam detik |

**Response `202 Accepted`:**
```json
{
  "success": true,
  "data": {
    "analysis_id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
    "status": "pending",
    "estimated_seconds": 8,
    "mode": "humming"
  }
}
```

**Error `400`** — Rekaman terlalu pendek:
```json
{
  "success": false,
  "error": {
    "code": "AUDIO_TOO_SHORT",
    "message": "Rekaman terlalu pendek. Minimal 5 detik diperlukan untuk analisis yang akurat",
    "details": {
      "received_seconds": 3,
      "minimum_seconds": 5
    }
  }
}
```

---

#### `GET /analyses/{analysis_id}`
Mengambil status dan hasil analisis. Digunakan untuk polling hingga proses selesai.

**Request:**
```http
GET /api/v1/analyses/018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b
X-Session-ID: 550e8400-e29b-41d4-a716-446655440000
```

**Response `200 OK`** — Status `pending` atau `processing`:
```json
{
  "success": true,
  "data": {
    "id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
    "status": "processing",
    "input_type": "youtube",
    "input_source": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "created_at": "2026-05-10T14:23:00Z"
  }
}
```

**Response `200 OK`** — Status `completed`:
```json
{
  "success": true,
  "data": {
    "id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
    "status": "completed",
    "input_type": "youtube",
    "input_source": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "result": {
      "detected_maqam": {
        "id": "hijaz",
        "name_arabic": "حجاز",
        "name_latin": "Hijaz",
        "name_indonesia": "Hijaz"
      },
      "confidence_score": 0.87,
      "confidence_label": "tinggi",
      "candidates": [
        { "rank": 1, "maqam_id": "hijaz",  "name_latin": "Hijaz",  "confidence_score": 0.87 },
        { "rank": 2, "maqam_id": "kurd",   "name_latin": "Kurd",   "confidence_score": 0.09 },
        { "rank": 3, "maqam_id": "bayati", "name_latin": "Bayati", "confidence_score": 0.04 }
      ],
      "explanation": {
        "karakteristik": "Maqam Hijaz dikenal dengan karakteristiknya yang dramatis dan agung...",
        "emosi": ["dramatis", "kerinduan", "agung", "spiritual"],
        "struktur_tangga_nada": "D – E♭ – F# – G – A – B♭ – C",
        "interval_khas": "Ciri khas Hijaz adalah interval augmented second antara E♭ dan F#...",
        "contoh_lagu": ["Ya Hanana", "Lir-ilir versi Arab", "Qasidah Burda segmen pembuka"],
        "tips_aransemen_banjari": "Maqam Hijaz sangat cocok untuk bagian pembuka yang dramatik..."
      },
      "audio_quality": "good",
      "processing_ms": 18240
    },
    "created_at": "2026-05-10T14:23:00Z",
    "completed_at": "2026-05-10T14:23:18Z"
  }
}
```

**Response `200 OK`** — Status `failed`:
```json
{
  "success": true,
  "data": {
    "id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
    "status": "failed",
    "error": {
      "code": "LOW_AUDIO_QUALITY",
      "message": "Kualitas audio terlalu rendah untuk deteksi maqam yang akurat."
    },
    "created_at": "2026-05-10T14:23:00Z"
  }
}
```

**Error `403`** — Analisis milik sesi lain:
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke analisis ini"
  }
}
```

> **Polling:** Frontend poll endpoint ini setiap 2 detik selama status `pending` atau `processing`. Hentikan jika status berubah ke `completed` atau `failed`, atau setelah timeout 90 detik.

---

### 2. Riwayat Analisis

#### `GET /history`
Mengambil daftar riwayat analisis milik sesi saat ini.

**Request:**
```http
GET /api/v1/history?page=1&limit=20&status=completed
X-Session-ID: 550e8400-e29b-41d4-a716-446655440000
```

**Query Parameters:**

| Parameter | Tipe | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | `integer` | `1` | Halaman pagination |
| `limit` | `integer` | `20` | Jumlah item per halaman (max: `50`) |
| `status` | `string` | — | Filter: `pending`, `processing`, `completed`, `failed` |
| `input_type` | `string` | — | Filter: `youtube`, `upload`, `microphone`, `humming` |

**Response `200 OK`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b",
      "status": "completed",
      "input_type": "youtube",
      "input_source": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
      "detected_maqam": {
        "id": "hijaz",
        "name_latin": "Hijaz",
        "name_indonesia": "Hijaz"
      },
      "confidence_score": 0.87,
      "created_at": "2026-05-10T14:23:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 7,
    "total_pages": 1
  }
}
```

---

#### `DELETE /history/{analysis_id}`
Menghapus satu entri riwayat analisis milik sesi saat ini.

**Request:**
```http
DELETE /api/v1/history/018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b
X-Session-ID: 550e8400-e29b-41d4-a716-446655440000
```

**Response `200 OK`:**
```json
{
  "success": true,
  "data": {
    "deleted_id": "018e1a2b-3c4d-7e5f-8a9b-0c1d2e3f4a5b"
  }
}
```

**Error `403`:**
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "Anda tidak memiliki akses ke riwayat ini"
  }
}
```

---

### 3. Referensi Maqam

#### `GET /maqamat`
Mengambil daftar semua maqam yang didukung sistem.

**Request:**
```http
GET /api/v1/maqamat
```

**Response `200 OK`:**
```json
{
  "success": true,
  "data": [
    {
      "id": "hijaz",
      "name_arabic": "حجاز",
      "name_latin": "Hijaz",
      "name_indonesia": "Hijaz",
      "emotion_tags": ["dramatis", "kerinduan", "agung"],
      "interval_description": "D – E♭ – F# – G – A – B♭ – C"
    },
    {
      "id": "rast",
      "name_arabic": "راست",
      "name_latin": "Rast",
      "name_indonesia": "Rast",
      "emotion_tags": ["tenang", "gembira", "natural"],
      "interval_description": "C – D – E♭(half) – F – G – A – B♭"
    }
  ],
  "meta": {
    "total": 8
  }
}
```

---

#### `GET /maqamat/{maqam_id}`
Mengambil detail lengkap satu maqam.

**Request:**
```http
GET /api/v1/maqamat/hijaz
```

**Response `200 OK`:**
```json
{
  "success": true,
  "data": {
    "id": "hijaz",
    "name_arabic": "حجاز",
    "name_latin": "Hijaz",
    "name_indonesia": "Hijaz",
    "interval_description": "D – E♭ – F# – G – A – B♭ – C",
    "characteristic_notes": ["E♭", "F#"],
    "emotion_tags": ["dramatis", "kerinduan", "agung", "spiritual"],
    "example_songs": [
      "Ya Hanana",
      "Tala'al Badru Alayna (versi Hijaz)",
      "Qasidah Burda"
    ],
    "tips_aransemen_banjari": "Maqam Hijaz sangat cocok untuk pembuka yang dramatis...",
    "pitch_class_profile": [1.0, 0.8, 0.0, 0.0, 0.6, 0.0, 0.9, 0.7, 0.0, 0.5, 0.0, 0.3,
                             0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0],
    "updated_at": "2026-05-01T00:00:00Z"
  }
}
```

**Error `404`:**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Maqam dengan ID 'unknown' tidak ditemukan"
  }
}
```

---

## Appendix A: Enum Values

| Field | Values |
|-------|--------|
| `input_type` | `youtube`, `upload`, `microphone`, `humming` |
| `status` (analyses) | `pending`, `processing`, `completed`, `failed` |
| `mode` (record) | `microphone`, `humming` |
| `confidence_label` | `sangat_rendah` (<40%), `rendah` (40–59%), `sedang` (60–74%), `tinggi` (75–89%), `sangat_tinggi` (≥90%) |
| `audio_quality` | `poor`, `fair`, `good`, `excellent` |

---

## Appendix B: Maqam IDs yang Didukung (v1)

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

## Appendix C: Rate Limiting

| Scope | Batas | Window |
|-------|-------|--------|
| Analisis per sesi | 10 request | Per jam |
| Upload file per sesi | 5 request | Per jam |

Response header saat mendekati limit:
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 2
X-RateLimit-Reset: 1715356800
```

---

## Appendix D: Polling Flow (Frontend)

```
POST /analyze/youtube → 202 { analysis_id, estimated_seconds }
         ↓
  tunggu 2 detik
         ↓
GET /analyses/{id} → status: "processing"
         ↓
  tunggu 2 detik
         ↓
GET /analyses/{id} → status: "completed" ✅
         ↓
  tampilkan hasil
```

Timeout polling: **90 detik**. Jika melebihi, tampilkan error dan sarankan user coba lagi.

---

**Versi History:**

| Versi | Tanggal | Perubahan |
|-------|---------|-----------|
| 1.1.0 | Mei 2026 | Hapus endpoint feedback & admin (F-013, bukan core). Hapus field `user_id` dari scope. Fokus murni core features F-001–F-007 |
| 1.0.0 | Mei 2026 | Draft awal |

*Versi 1.1.0 — Mei 2026*

"""Route analisis — endpoint internal yang dipanggil oleh Golang API."""

import os
import time
import traceback

from fastapi import APIRouter, File, Form, UploadFile, HTTPException

from app.services.pitch_extractor import extract_pcp_from_file
from app.services.maqam_matcher import match_maqam, get_confidence_label
from app.services.youtube_fetcher import (
    download_audio,
    cleanup_file,
    YouTubeFetchError,
)

router = APIRouter()

TEMP_DIR = os.environ.get("ANALYZER_TEMP_DIR", "/tmp/analyzer")
ALLOWED_AUDIO_MIMES = {
    "audio/mpeg",
    "audio/wav",
    "audio/x-wav",
    "audio/mp4",
    "audio/m4a",
    "audio/flac",
    "audio/ogg",
    "audio/webm",
    "application/octet-stream",  # Fallback browser recording
}


@router.post("/internal/analyze")
async def analyze_audio(
    source_type: str = Form(...),
    url: str = Form(None),
    file: UploadFile = File(None),
    segment_start: int = Form(0),
    segment_duration: int = Form(60),
    mode: str = Form("normal"),
):
    """
    Endpoint analisis internal — dipanggil oleh Golang API, bukan publik.

    Args:
        source_type: "youtube", "upload", "microphone", "humming"
        url: YouTube URL (wajib jika source_type=youtube)
        file: File audio upload/recording (wajib jika source_type!=youtube)
        segment_start: Detik mulai analisis
        segment_duration: Durasi segmen (max 120)
        mode: "normal", "microphone", atau "humming"
    """
    start_time = time.time()
    temp_file_path = None
    metadata = None

    try:
        # ── Step 1: Dapatkan file audio ──
        if source_type == "youtube":
            if not url:
                raise HTTPException(status_code=400, detail={
                    "code": "VALIDATION_ERROR",
                    "message": "URL YouTube wajib diisi untuk source_type=youtube",
                })
            try:
                result = download_audio(url)
                temp_file_path = result.file_path
                # Store metadata
                metadata = {
                    "title": result.title,
                    "duration": result.duration_seconds,
                    "channel": result.channel
                }
            except YouTubeFetchError as e:
                raise HTTPException(status_code=422, detail={
                    "code": e.code,
                    "message": e.message,
                })

        elif source_type in ("upload", "microphone", "humming"):
            if not file:
                raise HTTPException(status_code=400, detail={
                    "code": "VALIDATION_ERROR",
                    "message": "File audio wajib dikirim",
                })

            # Simpan ke temp file
            os.makedirs(TEMP_DIR, exist_ok=True)
            import uuid
            temp_file_path = os.path.join(TEMP_DIR, f"{uuid.uuid4()}.audio")

            content = await file.read()

            if len(content) > 50 * 1024 * 1024:  # 50MB
                raise HTTPException(status_code=413, detail={
                    "code": "FILE_TOO_LARGE",
                    "message": "Ukuran file melebihi batas 50MB",
                })

            with open(temp_file_path, "wb") as f:
                f.write(content)

            # Override mode berdasarkan source_type
            if source_type == "humming":
                mode = "humming"
            elif source_type == "microphone":
                mode = "microphone"
        else:
            raise HTTPException(status_code=400, detail={
                "code": "VALIDATION_ERROR",
                "message": f"source_type tidak valid: {source_type}",
            })

        # ── Step 2: Extract pitch → PCP ──
        segment_duration = min(segment_duration, 120)
        try:
            pcp, avg_confidence, audio_quality = extract_pcp_from_file(
                temp_file_path,
                segment_start=segment_start,
                segment_duration=segment_duration,
            )
        except ValueError as e:
            raise HTTPException(status_code=422, detail={
                "code": "LOW_AUDIO_QUALITY",
                "message": str(e),
            })

        # ── Step 3: Match maqam ──
        candidates = match_maqam(pcp, mode=mode, top_n=3)

        # ── Step 4: Format response ──
        processing_ms = int((time.time() - start_time) * 1000)

        top_candidate = candidates[0] if candidates else None
        confidence_label = (
            get_confidence_label(top_candidate.confidence_score)
            if top_candidate
            else "sangat_rendah"
        )

        return {
            "top3_candidates": [
                {
                    "maqam_id": c.maqam_id,
                    "name_latin": c.name_latin,
                    "name_arabic": c.name_arabic,
                    "confidence_score": c.confidence_score,
                    "rank": c.rank,
                }
                for c in candidates
            ],
            "confidence_label": confidence_label,
            "audio_quality": audio_quality,
            "processing_ms": processing_ms,
            "pcp": pcp.tolist(),
            "metadata": metadata if source_type == "youtube" else None,
        }

    except HTTPException:
        raise
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail={
            "code": "ANALYSIS_FAILED",
            "message": f"Terjadi kesalahan saat menganalisis: {str(e)}",
        })
    finally:
        # SELALU cleanup temp file
        if temp_file_path:
            cleanup_file(temp_file_path)

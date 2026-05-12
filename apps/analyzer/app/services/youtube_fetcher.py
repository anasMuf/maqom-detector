"""YouTube audio fetcher — download audio-only dari YouTube URL via yt-dlp."""

import os
import re
import uuid
import subprocess
import json
from dataclasses import dataclass

TEMP_DIR = os.environ.get("ANALYZER_TEMP_DIR", "/tmp/analyzer")
MAX_DURATION_SECONDS = 15 * 60  # 15 menit


@dataclass
class YouTubeAudioResult:
    """Hasil download audio dari YouTube."""

    file_path: str
    title: str
    duration_seconds: int
    channel: str


class YouTubeFetchError(Exception):
    """Base error untuk YouTube fetcher."""

    def __init__(self, message: str, code: str = "ANALYSIS_FAILED"):
        self.message = message
        self.code = code
        super().__init__(message)


def validate_youtube_url(url: str) -> bool:
    """Validasi apakah URL adalah YouTube yang valid."""
    patterns = [
        r"^https?://(www\.)?youtube\.com/watch\?v=[\w-]{11}",
        r"^https?://youtu\.be/[\w-]{11}",
        r"^https?://(www\.)?youtube\.com/shorts/[\w-]{11}",
        r"^https?://m\.youtube\.com/watch\?v=[\w-]{11}",
    ]
    return any(re.match(pattern, url) for pattern in patterns)


def get_video_info(url: str) -> dict:
    """Ambil metadata video tanpa download."""
    try:
        result = subprocess.run(
            [
                "yt-dlp",
                "--no-download",
                "--print-json",
                "--no-warnings",
                url,
            ],
            capture_output=True,
            text=True,
            timeout=30,
        )

        if result.returncode != 0:
            stderr = result.stderr.lower()
            if "private" in stderr or "sign in" in stderr:
                raise YouTubeFetchError(
                    "Video tidak dapat diakses. Kemungkinan video bersifat privat atau telah dihapus",
                    code="VIDEO_UNAVAILABLE",
                )
            if "not available" in stderr or "unavailable" in stderr:
                raise YouTubeFetchError(
                    "Video tidak tersedia",
                    code="VIDEO_UNAVAILABLE",
                )
            raise YouTubeFetchError(
                f"Gagal mengambil info video: {result.stderr[:200]}",
                code="ANALYSIS_FAILED",
            )

        return json.loads(result.stdout)

    except subprocess.TimeoutExpired:
        raise YouTubeFetchError(
            "Timeout saat mengambil info video",
            code="ANALYSIS_FAILED",
        )
    except json.JSONDecodeError:
        raise YouTubeFetchError(
            "Gagal membaca metadata video",
            code="ANALYSIS_FAILED",
        )


def download_audio(url: str) -> YouTubeAudioResult:
    """
    Download audio-only dari YouTube URL.

    Returns:
        YouTubeAudioResult dengan path ke file audio temp.

    Raises:
        YouTubeFetchError jika download gagal.
    """
    if not validate_youtube_url(url):
        raise YouTubeFetchError(
            "URL yang diberikan bukan URL YouTube yang valid",
            code="INVALID_URL",
        )

    # Ambil info video dulu
    info = get_video_info(url)
    duration = info.get("duration", 0)
    title = info.get("title", "Unknown")
    channel = info.get("channel", info.get("uploader", "Unknown"))

    if duration > MAX_DURATION_SECONDS:
        raise YouTubeFetchError(
            f"Durasi video terlalu panjang ({duration // 60} menit). Maksimal 15 menit.",
            code="VIDEO_TOO_LONG",
        )

    # Buat temp directory
    os.makedirs(TEMP_DIR, exist_ok=True)
    file_id = str(uuid.uuid4())
    output_path = os.path.join(TEMP_DIR, f"{file_id}.wav")

    try:
        result = subprocess.run(
            [
                "yt-dlp",
                "-x",  # extract audio
                "--audio-format",
                "wav",
                "--audio-quality",
                "0",
                "-o",
                output_path,
                "--no-playlist",
                "--no-warnings",
                url,
            ],
            capture_output=True,
            text=True,
            timeout=120,
        )

        if result.returncode != 0:
            raise YouTubeFetchError(
                f"Gagal mengunduh audio: {result.stderr[:200]}",
                code="ANALYSIS_FAILED",
            )

        # yt-dlp mungkin menambahkan ekstensi — cari file yang sebenarnya
        actual_path = output_path
        if not os.path.exists(actual_path):
            # Cari file dengan UUID yang sama
            for f in os.listdir(TEMP_DIR):
                if f.startswith(file_id):
                    actual_path = os.path.join(TEMP_DIR, f)
                    break

        if not os.path.exists(actual_path):
            raise YouTubeFetchError(
                "File audio tidak ditemukan setelah download",
                code="ANALYSIS_FAILED",
            )

        return YouTubeAudioResult(
            file_path=actual_path,
            title=title,
            duration_seconds=duration,
            channel=channel,
        )

    except subprocess.TimeoutExpired:
        cleanup_file(output_path)
        raise YouTubeFetchError(
            "Timeout saat mengunduh audio (melebihi 2 menit)",
            code="ANALYSIS_FAILED",
        )


def cleanup_file(file_path: str) -> None:
    """Hapus file temp. Selalu aman untuk dipanggil meskipun file tidak ada."""
    try:
        if file_path and os.path.exists(file_path):
            os.remove(file_path)
    except OSError:
        pass

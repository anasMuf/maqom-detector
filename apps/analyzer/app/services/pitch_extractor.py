"""Pitch extractor — extract pitch dari audio menggunakan CREPE + librosa."""

import numpy as np
import librosa
import crepe
import os
from pydub import AudioSegment

# Konfigurasi default
DEFAULT_SR = 22050  # Sample rate standar librosa
CREPE_MODEL = "tiny"  # Model CREPE paling ringan (~20MB RAM)
CREPE_CONFIDENCE_THRESHOLD = 0.50  # Filter frame dengan confidence rendah
NUM_PITCH_CLASSES = 12  # 12-TET (C, C#, D, ..., B)


def load_audio(
    file_path: str,
    segment_start: int = 0,
    segment_duration: int = 60,
) -> np.ndarray:
    """
    Load audio file dan crop ke segmen yang diminta.
    Mendukung auto-conversion dari MP3/M4A ke WAV via pydub.
    """
    ext = os.path.splitext(file_path)[1].lower()
    temp_wav = None
    target_path = file_path

    # Jika bukan WAV, konversi ke WAV temp menggunakan pydub
    if ext != ".wav" and ext != "":
        try:
            audio_segment = AudioSegment.from_file(file_path)
            temp_wav = file_path + ".temp.wav"
            audio_segment.export(temp_wav, format="wav")
            target_path = temp_wav
        except Exception as e:
            print(f"Warning: Gagal konversi {ext} ke WAV via pydub: {e}")

    try:
        # Load audio mono
        y, sr = librosa.load(
            target_path,
            sr=DEFAULT_SR,
            mono=True,
            offset=segment_start,
            duration=segment_duration,
        )

        if len(y) < DEFAULT_SR * 3:
            raise ValueError(
                "Audio terlalu pendek setelah cropping. Minimal 3 detik diperlukan."
            )

        return y
    finally:
        # Cleanup temp file jika ada
        if temp_wav and os.path.exists(temp_wav):
            try:
                os.remove(temp_wav)
            except:
                pass


def extract_pitch(audio: np.ndarray) -> tuple[np.ndarray, np.ndarray, np.ndarray]:
    """
    Extract pitch menggunakan CREPE.

    Args:
        audio: Audio signal (mono, sr=22050)

    Returns:
        Tuple (time, frequency, confidence) arrays
    """
    # CREPE butuh sr=16000 untuk model tiny
    audio_16k = librosa.resample(audio, orig_sr=DEFAULT_SR, target_sr=16000)

    time, frequency, confidence, _ = crepe.predict(
        audio_16k,
        sr=16000,
        model_capacity=CREPE_MODEL,
        viterbi=True,  # Smoothing untuk hasil lebih stabil
        step_size=10,  # 10ms per frame
    )

    return time, frequency, confidence


def hz_to_pitch_class(frequency: float) -> int:
    """
    Konversi frekuensi (Hz) ke pitch class (0-11).

    0=C, 1=C#, 2=D, ..., 11=B
    """
    if frequency <= 0:
        return -1
    midi = 69 + 12 * np.log2(frequency / 440.0)
    return int(round(midi)) % 12


def compute_pitch_class_profile(
    frequency: np.ndarray,
    confidence: np.ndarray,
    confidence_threshold: float = CREPE_CONFIDENCE_THRESHOLD,
) -> tuple[np.ndarray, float, str]:
    """
    Hitung Pitch Class Profile (PCP) dari hasil CREPE.

    PCP adalah histogram 12-bin yang menunjukkan distribusi pitch classes
    dalam audio. Setiap bin mewakili satu semitone (C, C#, D, ..., B).

    Args:
        frequency: Array frekuensi per frame (Hz)
        confidence: Array confidence per frame (0.0–1.0)
        confidence_threshold: Threshold minimum confidence frame

    Returns:
        Tuple (pcp_normalized, avg_confidence, audio_quality)
    """
    # Filter frame dengan confidence rendah
    mask = confidence >= confidence_threshold
    valid_freq = frequency[mask]
    valid_conf = confidence[mask]

    if len(valid_freq) < 10:
        raise ValueError(
            "Terlalu sedikit frame pitch yang valid. "
            "Audio mungkin terlalu noisy atau tidak mengandung melodi."
        )

    # Hitung PCP — weighted by confidence
    pcp = np.zeros(NUM_PITCH_CLASSES, dtype=np.float64)

    for freq, conf in zip(valid_freq, valid_conf):
        pc = hz_to_pitch_class(freq)
        if pc >= 0:
            pcp[pc] += conf  # Weight by confidence

    # Normalize ke range 0.0–1.0
    pcp_max = pcp.max()
    if pcp_max > 0:
        pcp_normalized = pcp / pcp_max
    else:
        pcp_normalized = pcp

    # Hitung audio quality dari rata-rata confidence
    avg_confidence = float(np.mean(valid_conf))
    audio_quality = _classify_audio_quality(avg_confidence, len(valid_freq), len(frequency))

    return pcp_normalized, avg_confidence, audio_quality


def _classify_audio_quality(
    avg_confidence: float,
    valid_frames: int,
    total_frames: int,
) -> str:
    """Klasifikasi kualitas audio berdasarkan CREPE confidence."""
    valid_ratio = valid_frames / max(total_frames, 1)

    if avg_confidence >= 0.85 and valid_ratio >= 0.7:
        return "excellent"
    elif avg_confidence >= 0.70 and valid_ratio >= 0.5:
        return "good"
    elif avg_confidence >= 0.55 and valid_ratio >= 0.3:
        return "fair"
    else:
        return "poor"


def extract_pcp_from_file(
    file_path: str,
    segment_start: int = 0,
    segment_duration: int = 60,
) -> tuple[np.ndarray, float, str]:
    """
    Pipeline lengkap: file audio → PCP.

    Args:
        file_path: Path ke file audio
        segment_start: Detik mulai analisis
        segment_duration: Durasi segmen dalam detik

    Returns:
        Tuple (pcp, avg_confidence, audio_quality)
    """
    audio = load_audio(file_path, segment_start, segment_duration)
    _, frequency, confidence = extract_pitch(audio)
    return compute_pitch_class_profile(frequency, confidence)

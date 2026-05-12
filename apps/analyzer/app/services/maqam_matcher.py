"""Maqam matcher — cocokkan PCP input vs template maqam menggunakan cosine similarity."""

import json
import os
from dataclasses import dataclass

import numpy as np

TEMPLATES_PATH = os.path.join(
    os.path.dirname(__file__), "..", "data", "maqam_templates.json"
)

# Threshold matching
HUMMING_SIMILARITY_DISCOUNT = 0.85  # Toleransi lebih longgar untuk mode humming


@dataclass
class MaqamCandidate:
    """Satu kandidat maqam hasil matching."""

    maqam_id: str
    name_latin: str
    name_arabic: str
    confidence_score: float
    rank: int
    best_transposition: int  # Semitone offset yang menghasilkan similarity tertinggi


def load_templates() -> list[dict]:
    """Load maqam templates dari JSON file."""
    with open(TEMPLATES_PATH, "r", encoding="utf-8") as f:
        return json.load(f)


def cosine_similarity(a: np.ndarray, b: np.ndarray) -> float:
    """Hitung cosine similarity antara dua vector."""
    dot = np.dot(a, b)
    norm_a = np.linalg.norm(a)
    norm_b = np.linalg.norm(b)

    if norm_a == 0 or norm_b == 0:
        return 0.0

    return float(dot / (norm_a * norm_b))


def rotate_pcp(pcp: np.ndarray, semitones: int) -> np.ndarray:
    """
    Rotasi PCP sebanyak N semitone.

    Ini mensimulasikan transposisi: jika lagu dimainkan di key berbeda,
    PCP-nya akan ter-shift secara sirkuler.
    """
    return np.roll(pcp, -semitones)


def match_maqam(
    input_pcp: np.ndarray,
    mode: str = "normal",
    top_n: int = 3,
) -> list[MaqamCandidate]:
    """
    Cocokkan PCP input dengan semua template maqam.

    Untuk setiap maqam template, coba semua 12 rotasi (transposisi)
    dan ambil rotasi dengan similarity tertinggi. Ini memungkinkan
    deteksi maqam terlepas dari key/tonic yang digunakan.

    Args:
        input_pcp: Pitch Class Profile 12-bin dari audio input
        mode: "normal", "microphone", atau "humming"
        top_n: Jumlah kandidat teratas yang dikembalikan

    Returns:
        List MaqamCandidate terurut descending berdasarkan confidence
    """
    templates = load_templates()
    results: list[tuple[float, int, dict]] = []

    for template in templates:
        template_pcp = np.array(template["pitch_class_profile"], dtype=np.float64)

        best_similarity = -1.0
        best_transposition = 0

        # Coba semua 12 rotasi (transposisi)
        for shift in range(12):
            rotated_input = rotate_pcp(input_pcp, shift)
            sim = cosine_similarity(rotated_input, template_pcp)

            if sim > best_similarity:
                best_similarity = sim
                best_transposition = shift

        results.append((best_similarity, best_transposition, template))

    # Sort descending by similarity
    results.sort(key=lambda x: x[0], reverse=True)

    # Normalisasi confidence scores
    # Gunakan softmax-like normalization agar top candidates lebih terdiferensiasi
    similarities = np.array([r[0] for r in results])
    total_sim = similarities.sum()

    candidates = []
    for rank, (similarity, transposition, template) in enumerate(results[:top_n], 1):
        # Confidence = proporsi relatif terhadap total similarity
        confidence = float(similarity / total_sim) if total_sim > 0 else 0.0

        # Untuk mode humming, diskon confidence (lebih uncertain)
        if mode == "humming":
            confidence *= HUMMING_SIMILARITY_DISCOUNT

        candidates.append(
            MaqamCandidate(
                maqam_id=template["id"],
                name_latin=template["name_latin"],
                name_arabic=template["name_arabic"],
                confidence_score=round(confidence, 4),
                rank=rank,
                best_transposition=transposition,
            )
        )

    return candidates


def get_confidence_label(confidence: float) -> str:
    """Konversi confidence score ke label verbal."""
    if confidence >= 0.90:
        return "sangat_tinggi"
    elif confidence >= 0.75:
        return "tinggi"
    elif confidence >= 0.60:
        return "sedang"
    elif confidence >= 0.40:
        return "rendah"
    else:
        return "sangat_rendah"

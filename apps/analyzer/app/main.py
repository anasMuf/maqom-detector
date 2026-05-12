"""MaqamDetector — Python Analyzer Service"""

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.routes.analyze import router as analyze_router

app = FastAPI(
    title="MaqamDetector Analyzer",
    description="Internal pitch extraction and maqam matching service",
    version="1.0.0",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Register routes
app.include_router(analyze_router)


@app.get("/health")
async def health_check():
    """Health check endpoint for Docker healthcheck and readiness probes."""
    return {"status": "ok", "service": "analyzer"}

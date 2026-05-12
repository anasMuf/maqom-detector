package handler

import (
	"api/dto"
	"api/middleware"
	"api/service"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AnalyzeHandler struct {
	analyzeService *service.AnalyzeService
}

func NewAnalyzeHandler(analyzeService *service.AnalyzeService) *AnalyzeHandler {
	return &AnalyzeHandler{analyzeService: analyzeService}
}

// AnalyzeYoutube godoc
// @Summary     Analisis maqam dari YouTube URL
// @Description Memulai proses analisis maqam dari video YouTube secara asinkron
// @Tags        analyze
// @Accept      json
// @Produce     json
// @Param       X-Session-ID header string true "Session UUID"
// @Param       body body dto.AnalyzeYoutubeRequest true "YouTube analysis request"
// @Success     202 {object} dto.AnalyzeAcceptedResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     429 {object} dto.ErrorResponse
// @Router      /analyze/youtube [post]
func (h *AnalyzeHandler) AnalyzeYoutube(c echo.Context) error {
	sessionID := middleware.GetSessionID(c)

	var req dto.AnalyzeYoutubeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "VALIDATION_ERROR", "message": "Request body tidak valid"},
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "VALIDATION_ERROR", "message": err.Error()},
		})
	}

	if req.SegmentDuration == 0 {
		req.SegmentDuration = 60
	}

	analysis, err := h.analyzeService.StartYoutubeAnalysis(sessionID, req.URL, req.SegmentStart, req.SegmentDuration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "ANALYSIS_FAILED", "message": err.Error()},
		})
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"success": true,
		"data": dto.AnalyzeAcceptedResponse{
			AnalysisID:       analysis.ID.String(),
			Status:           "pending",
			EstimatedSeconds: 30,
			PollingURL:       fmt.Sprintf("/api/v1/analyses/%s", analysis.ID),
		},
	})
}

// AnalyzeUpload godoc
// @Summary     Analisis maqam dari file audio upload
// @Tags        analyze
// @Accept      multipart/form-data
// @Produce     json
// @Param       X-Session-ID header string true "Session UUID"
// @Param       file formData file true "File audio (WAV, MP3, M4A, dll.)"
// @Success     202 {object} dto.AnalyzeAcceptedResponse
// @Failure     400 {object} dto.ErrorResponse
// @Failure     413 {object} dto.ErrorResponse
// @Router      /analyze/upload [post]
func (h *AnalyzeHandler) AnalyzeUpload(c echo.Context) error {
	return h.handleFileAnalysis(c, "upload")
}

// AnalyzeRecord godoc
// @Summary     Analisis maqam dari rekaman mikrofon/humming
// @Tags        analyze
// @Accept      multipart/form-data
// @Produce     json
// @Param       X-Session-ID header string true "Session UUID"
// @Param       file formData file true "Audio recording"
// @Param       mode formData string false "Mode: microphone atau humming" default(microphone)
// @Success     202 {object} dto.AnalyzeAcceptedResponse
// @Failure     400 {object} dto.ErrorResponse
// @Router      /analyze/record [post]
func (h *AnalyzeHandler) AnalyzeRecord(c echo.Context) error {
	mode := c.FormValue("mode")
	if mode == "humming" {
		return h.handleFileAnalysis(c, "humming")
	}
	return h.handleFileAnalysis(c, "microphone")
}

func (h *AnalyzeHandler) handleFileAnalysis(c echo.Context, sourceType string) error {
	sessionID := middleware.GetSessionID(c)

	file, err := c.FormFile("file")
	if err != nil || file == nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "VALIDATION_ERROR", "message": "File audio wajib dikirim"},
		})
	}

	// Check file size (50MB max)
	if file.Size > 50*1024*1024 {
		return c.JSON(http.StatusRequestEntityTooLarge, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "FILE_TOO_LARGE", "message": "Ukuran file melebihi batas 50MB"},
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "ANALYSIS_FAILED", "message": "Gagal membaca file"},
		})
	}
	defer src.Close()

	fileContent, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "ANALYSIS_FAILED", "message": "Gagal membaca konten file"},
		})
	}

	analysis, err := h.analyzeService.StartUploadAnalysis(sessionID, file.Filename, fileContent, sourceType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "ANALYSIS_FAILED", "message": err.Error()},
		})
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"success": true,
		"data": dto.AnalyzeAcceptedResponse{
			AnalysisID:       analysis.ID.String(),
			Status:           "pending",
			EstimatedSeconds: 30,
			PollingURL:       fmt.Sprintf("/api/v1/analyses/%s", analysis.ID),
		},
	})
}

// GetAnalysis godoc
// @Summary     Ambil status/hasil analisis (polling)
// @Tags        analyze
// @Produce     json
// @Param       X-Session-ID header string true "Session UUID"
// @Param       id path string true "Analysis UUID"
// @Success     200 {object} dto.AnalysisDetailResponse
// @Failure     404 {object} dto.ErrorResponse
// @Router      /analyses/{id} [get]
func (h *AnalyzeHandler) GetAnalysis(c echo.Context) error {
	sessionID := middleware.GetSessionID(c)

	analysisID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "VALIDATION_ERROR", "message": "ID analisis tidak valid"},
		})
	}

	analysis, err := h.analyzeService.GetAnalysis(analysisID, sessionID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "NOT_FOUND", "message": "Analisis tidak ditemukan"},
		})
	}

	// Build response
	resp := dto.AnalysisDetailResponse{
		ID:              analysis.ID.String(),
		Status:          analysis.Status,
		InputType:       analysis.InputType,
		InputSource:     analysis.InputSource,
		DetectedMaqamID: analysis.DetectedMaqamID,
		ConfidenceScore: analysis.ConfidenceScore,
		ConfidenceLabel: analysis.ConfidenceLabel,
		ExplanationText: analysis.ExplanationText,
		AudioQuality:    analysis.AudioQuality,
		ProcessingMs:    analysis.ProcessingMs,
		ErrorCode:       analysis.ErrorCode,
		ErrorMessage:    analysis.ErrorMessage,
		CreatedAt:       analysis.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if analysis.CompletedAt != nil {
		t := analysis.CompletedAt.Format("2006-01-02T15:04:05Z")
		resp.CompletedAt = &t
	}

	for _, c := range analysis.Candidates {
		resp.Candidates = append(resp.Candidates, dto.CandidateResponse{
			MaqamID:         c.MaqamID,
			NameLatin:       c.Maqam.NameLatin,
			NameArabic:      c.Maqam.NameArabic,
			ConfidenceScore: c.ConfidenceScore,
			Rank:            c.Rank,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

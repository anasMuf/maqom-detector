package handler

import (
	"api/dto"
	"api/middleware"
	"api/service"
	"math"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type HistoryHandler struct {
	historyService *service.HistoryService
}

func NewHistoryHandler(historyService *service.HistoryService) *HistoryHandler {
	return &HistoryHandler{historyService: historyService}
}

// GetHistory godoc
// @Summary     Riwayat analisis
// @Tags        history
// @Produce     json
// @Param       X-Session-ID header string true "Session UUID"
// @Param       page query int false "Halaman" default(1)
// @Param       limit query int false "Jumlah per halaman" default(10)
// @Success     200 {object} dto.HistoryListResponse
// @Router      /history [get]
func (h *HistoryHandler) GetHistory(c echo.Context) error {
	sessionID := middleware.GetSessionID(c)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	analyses, total, err := h.historyService.GetHistory(sessionID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "ANALYSIS_FAILED", "message": "Gagal mengambil riwayat"},
		})
	}

	var items []dto.HistoryItemResponse
	for _, a := range analyses {
		item := dto.HistoryItemResponse{
			ID:              a.ID.String(),
			InputType:       a.InputType,
			InputSource:     a.InputSource,
			DetectedMaqamID: a.DetectedMaqamID,
			ConfidenceScore: a.ConfidenceScore,
			Status:          a.Status,
			CreatedAt:       a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}

		if a.DetectedMaqam != nil {
			item.MaqamNameLatin = a.DetectedMaqam.NameLatin
		}

		items = append(items, item)
	}

	if items == nil {
		items = []dto.HistoryItemResponse{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"items": items,
			"meta": dto.PaginationMeta{
				Page:       page,
				Limit:      limit,
				TotalItems: total,
				TotalPages: totalPages,
			},
		},
	})
}

// DeleteHistory godoc
// @Summary     Hapus riwayat analisis
// @Tags        history
// @Produce     json
// @Param       X-Session-ID header string true "Session UUID"
// @Param       id path string true "Analysis UUID"
// @Success     200 {object} dto.SuccessResponse
// @Failure     404 {object} dto.ErrorResponse
// @Router      /history/{id} [delete]
func (h *HistoryHandler) DeleteHistory(c echo.Context) error {
	sessionID := middleware.GetSessionID(c)

	analysisID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "VALIDATION_ERROR", "message": "ID tidak valid"},
		})
	}

	if err := h.historyService.DeleteHistory(sessionID, analysisID); err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "NOT_FOUND", "message": "Riwayat analisis tidak ditemukan"},
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    map[string]string{"message": "Riwayat berhasil dihapus"},
	})
}

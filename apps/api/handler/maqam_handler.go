package handler

import (
	"api/dto"
	"api/service"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type MaqamHandler struct {
	maqamService *service.MaqamService
}

func NewMaqamHandler(maqamService *service.MaqamService) *MaqamHandler {
	return &MaqamHandler{maqamService: maqamService}
}

// GetMaqamat godoc
// @Summary     Daftar semua maqam yang didukung
// @Tags        maqam
// @Produce     json
// @Success     200 {array} dto.MaqamListResponse
// @Router      /maqamat [get]
func (h *MaqamHandler) GetMaqamat(c echo.Context) error {
	maqams, err := h.maqamService.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "ANALYSIS_FAILED", "message": "Gagal mengambil data maqam"},
		})
	}

	var items []dto.MaqamListResponse
	for _, m := range maqams {
		var emotionTags []string
		json.Unmarshal([]byte(m.EmotionTags), &emotionTags)

		items = append(items, dto.MaqamListResponse{
			ID:                  m.ID,
			NameLatin:           m.NameLatin,
			NameArabic:          m.NameArabic,
			NameIndonesia:       m.NameIndonesia,
			IntervalDescription: m.IntervalDescription,
			EmotionTags:         emotionTags,
		})
	}

	if items == nil {
		items = []dto.MaqamListResponse{}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    items,
	})
}

// GetMaqamByID godoc
// @Summary     Detail maqam
// @Tags        maqam
// @Produce     json
// @Param       id path string true "Maqam ID (e.g. hijaz)"
// @Success     200 {object} dto.MaqamDetailResponse
// @Failure     404 {object} dto.ErrorResponse
// @Router      /maqamat/{id} [get]
func (h *MaqamHandler) GetMaqamByID(c echo.Context) error {
	id := c.Param("id")

	maqam, err := h.maqamService.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"error":   map[string]string{"code": "NOT_FOUND", "message": "Maqam tidak ditemukan"},
		})
	}

	var emotionTags []string
	json.Unmarshal([]byte(maqam.EmotionTags), &emotionTags)

	var exampleSongs []string
	json.Unmarshal([]byte(maqam.ExampleSongs), &exampleSongs)

	resp := dto.MaqamDetailResponse{
		ID:                  maqam.ID,
		NameLatin:           maqam.NameLatin,
		NameArabic:          maqam.NameArabic,
		NameIndonesia:       maqam.NameIndonesia,
		IntervalDescription: maqam.IntervalDescription,
		EmotionTags:         emotionTags,
		ExampleSongs:        exampleSongs,
		TipsAransemen:       maqam.TipsAransemen,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

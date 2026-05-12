package repository

import (
	"api/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnalysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Create(analysis *model.Analysis) error {
	return r.db.Create(analysis).Error
}

func (r *AnalysisRepository) FindByID(id uuid.UUID) (*model.Analysis, error) {
	analysis := &model.Analysis{}
	err := r.db.Preload("DetectedMaqam").Preload("Candidates.Maqam").
		Where("id = ?", id).First(analysis).Error
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

func (r *AnalysisRepository) FindByIDAndSession(id, sessionID uuid.UUID) (*model.Analysis, error) {
	analysis := &model.Analysis{}
	err := r.db.Preload("DetectedMaqam").Preload("Candidates.Maqam").
		Where("id = ? AND session_id = ?", id, sessionID).First(analysis).Error
	if err != nil {
		return nil, err
	}
	return analysis, nil
}

func (r *AnalysisRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Analysis{}).Where("id = ?", id).
		Update("status", status).Error
}

// AnalysisResult used to update analysis with final results
type AnalysisResult struct {
	DetectedMaqamID string
	ConfidenceScore float64
	ConfidenceLabel string
	ExplanationText string
	AudioQuality    string
	ProcessingMs    int
	Candidates      []model.AnalysisCandidate
}

func (r *AnalysisRepository) UpdateCompleted(id uuid.UUID, result *AnalysisResult) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := gorm.Expr("NOW()")

		if err := tx.Model(&model.Analysis{}).Where("id = ?", id).Updates(map[string]interface{}{
			"detected_maqam_id": result.DetectedMaqamID,
			"confidence_score":  result.ConfidenceScore,
			"confidence_label":  result.ConfidenceLabel,
			"explanation_text":  result.ExplanationText,
			"audio_quality":     result.AudioQuality,
			"processing_ms":     result.ProcessingMs,
			"status":            "completed",
			"completed_at":      now,
		}).Error; err != nil {
			return err
		}

		// Save candidates
		for i := range result.Candidates {
			result.Candidates[i].AnalysisID = id
			if err := tx.Create(&result.Candidates[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *AnalysisRepository) UpdateFailed(id uuid.UUID, errorCode, errorMessage string) error {
	now := gorm.Expr("NOW()")
	return r.db.Model(&model.Analysis{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        "failed",
		"error_code":    errorCode,
		"error_message": errorMessage,
		"completed_at":  now,
	}).Error
}

// HistoryFilter for paginated history queries
type HistoryFilter struct {
	Page  int
	Limit int
}

func (r *AnalysisRepository) FindBySession(sessionID uuid.UUID, filter HistoryFilter) ([]model.Analysis, int64, error) {
	var analyses []model.Analysis
	var total int64

	query := r.db.Model(&model.Analysis{}).Where("session_id = ?", sessionID)
	query.Count(&total)

	offset := (filter.Page - 1) * filter.Limit
	err := query.Preload("DetectedMaqam").Preload("Candidates.Maqam").
		Order("created_at DESC").
		Offset(offset).Limit(filter.Limit).
		Find(&analyses).Error

	return analyses, total, err
}

func (r *AnalysisRepository) DeleteByIDAndSession(id, sessionID uuid.UUID) error {
	result := r.db.Where("id = ? AND session_id = ?", id, sessionID).
		Delete(&model.Analysis{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

package service

import (
	"api/model"
	"api/repository"

	"github.com/google/uuid"
)

type HistoryService struct {
	analysisRepo *repository.AnalysisRepository
}

func NewHistoryService(analysisRepo *repository.AnalysisRepository) *HistoryService {
	return &HistoryService{analysisRepo: analysisRepo}
}

func (s *HistoryService) GetHistory(sessionID uuid.UUID, page, limit int) ([]model.Analysis, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	filter := repository.HistoryFilter{
		Page:  page,
		Limit: limit,
	}

	return s.analysisRepo.FindBySession(sessionID, filter)
}

func (s *HistoryService) DeleteHistory(sessionID, analysisID uuid.UUID) error {
	return s.analysisRepo.DeleteByIDAndSession(analysisID, sessionID)
}

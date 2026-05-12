package model

import "github.com/google/uuid"

// AnalysisCandidate — top-N kandidat maqam dari satu analisis
type AnalysisCandidate struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AnalysisID      uuid.UUID `json:"analysis_id" gorm:"type:uuid;not null;index"`
	MaqamID         string    `json:"maqam_id" gorm:"type:varchar(50);not null"`
	Maqam           Maqam     `json:"maqam" gorm:"foreignKey:MaqamID"`
	ConfidenceScore float64   `json:"confidence_score"`
	Rank            int       `json:"rank"` // 1, 2, 3
}

func (AnalysisCandidate) TableName() string {
	return "analysis_candidates"
}

package model

import (
	"time"

	"github.com/google/uuid"
)

// Analysis — record analisis maqam
type Analysis struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SessionID       uuid.UUID  `json:"session_id" gorm:"type:uuid;not null;index"`
	Session         Session    `json:"-" gorm:"foreignKey:SessionID"`
	InputType       string     `json:"input_type" gorm:"type:varchar(20);not null"`       // youtube|upload|microphone|humming
	InputSource     string     `json:"input_source" gorm:"type:text"`                     // URL atau filename
	DetectedMaqamID *string    `json:"detected_maqam_id" gorm:"type:varchar(50)"`
	DetectedMaqam   *Maqam     `json:"detected_maqam,omitempty" gorm:"foreignKey:DetectedMaqamID"`
	ConfidenceScore *float64   `json:"confidence_score"`
	ConfidenceLabel string     `json:"confidence_label" gorm:"type:varchar(20)"`          // sangat_tinggi|tinggi|sedang|rendah|sangat_rendah
	ExplanationText string     `json:"explanation_text" gorm:"type:text"`
	AudioQuality    string     `json:"audio_quality" gorm:"type:varchar(20)"`             // excellent|good|fair|poor
	Status          string     `json:"status" gorm:"type:varchar(20);not null;default:'pending'"` // pending|processing|completed|failed
	ProcessingMs    *int       `json:"processing_ms"`
	ErrorCode       string     `json:"error_code,omitempty" gorm:"type:varchar(50)"`
	ErrorMessage    string     `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`

	// Relations
	Candidates []AnalysisCandidate `json:"candidates,omitempty" gorm:"foreignKey:AnalysisID"`
}

func (Analysis) TableName() string {
	return "analyses"
}

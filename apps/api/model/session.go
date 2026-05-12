package model

import (
	"time"

	"github.com/google/uuid"
)

// Session — guest session identifikasi via X-Session-ID header
type Session struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LastActiveAt time.Time  `json:"last_active_at" gorm:"autoUpdateTime"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

func (Session) TableName() string {
	return "sessions"
}

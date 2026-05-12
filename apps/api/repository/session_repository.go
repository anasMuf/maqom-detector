package repository

import (
	"api/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// FindOrCreate finds an existing session or creates a new one
func (r *SessionRepository) FindOrCreate(sessionID uuid.UUID) (*model.Session, error) {
	session := &model.Session{}
	result := r.db.Where("id = ?", sessionID).First(session)

	if result.Error == gorm.ErrRecordNotFound {
		session.ID = sessionID
		session.LastActiveAt = time.Now()
		if err := r.db.Create(session).Error; err != nil {
			return nil, err
		}
		return session, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	// Update last active
	r.db.Model(session).Update("last_active_at", time.Now())
	return session, nil
}

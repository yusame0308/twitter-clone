package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tweet struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tweet) BeforeCreate(tx *gorm.DB) error {
	newUUID := uuid.New()
	t.ID = newUUID.String()
	return nil
}

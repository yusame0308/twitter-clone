package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Tweet struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *Tweet) BeforeCreate(*gorm.DB) error {
	newUUID := uuid.New()
	t.ID = newUUID.String()
	return nil
}

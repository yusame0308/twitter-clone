package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	newUUID := uuid.New()
	u.ID = newUUID.String()
	return nil
}

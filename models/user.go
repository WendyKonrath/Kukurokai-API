package models

import (
	"time"
	"github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
	Role      string    `json:"role" gorm:"not null" validate:"required,oneof=superadmin admin user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Gerar ID automaticamente com nanoid
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID, err = gonanoid.New()
	}
	return
}
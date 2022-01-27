package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserID struct {
	ID string `uri:"id" binding:"required"`
}

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	FirstName string         `json:"firstname" binding:"required"`
	LastName  string         `json:"lastname" binding:"required"`
	CreatedAt int64          `gorm:"autoCreateTime:milli" json:"created_at"`
	UpdatedAt int64          `gorm:"autoUpdateTime:milli" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (x *User) FillDefaults() {
	if x.ID == "" {
		x.ID = uuid.New().String()
	}
}

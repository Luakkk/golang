package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	CreatedAt time.Time      `json:"CreatedAt"`
	UpdatedAt time.Time      `json:"UpdatedAt"`
	DeletedAt gorm.DeletedAt `json:"DeletedAt" gorm:"index"`
	ID        uuid.UUID      `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string         `json:"Username"`
	Email     string         `json:"Email"`
	Password  string         `json:"Password"`
	Role      string         `json:"Role"` // user, admin, etc.
	Verified  bool           `json:"Verified"`
}

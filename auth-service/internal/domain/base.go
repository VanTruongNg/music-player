package domain

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

// BaseModel defines common fields for all domain models.
type BaseModel struct {
	ID        string      `json:"id" gorm:"type:char(26);primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// NewBaseModel creates a new BaseModel with generated ULID and current timestamps.
func NewBaseModel() BaseModel {
	return BaseModel{
		ID:        ulid.Make().String(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

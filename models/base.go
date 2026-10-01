package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID string `gorm:"primaryKey;size:32;not null" json:"id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID != "" {
		return nil
	}

	u, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to create uuid: %w", err)
	}
	b.ID = strings.ReplaceAll(u.String(), "-", "")

	return nil
}

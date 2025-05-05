package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        int            `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

var ErrUnsupportedFileType = errors.New("unsupported file type")

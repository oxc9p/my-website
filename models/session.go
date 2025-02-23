package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID            uint      `gorm:"primaryKey"`
	SessionID     string    `gorm:"uniqueIndex"`
	UserID        uuid.UUID `gorm:"not null"`
	UserName      string    `gorm:"not null"`
	Authenticated bool      `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}

func (s *Session) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return nil
}

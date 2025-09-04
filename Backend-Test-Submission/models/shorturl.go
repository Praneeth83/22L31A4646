package models

import "time"

type ShortURL struct {
	ID           uint   `gorm:"primaryKey"`
	OriginalURL  string `gorm:"not null"`
	ShortCode    string `gorm:"uniqueIndex;not null"`
	CreatedAt    time.Time
	ExpiryDate   time.Time
	Clicks       int
	ClickDetails []ClickInfo `gorm:"foreignKey:ShortURLID"`
}

package models

import "time"

type ClickInfo struct {
	ID         uint `gorm:"primaryKey"`
	ShortURLID uint
	Timestamp  time.Time
	SourceIP   string
	Location   string
}

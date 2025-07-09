package models

import "time"

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"size:512;uniqueIndex"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

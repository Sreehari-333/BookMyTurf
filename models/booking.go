package models

import (
	"time"

	"gorm.io/gorm"
)

// Creating struct for Booking details

type Booking struct {
	gorm.Model
	UserId    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey=UserId" json:"-" binding:"-"`
	TurfId    uint      `json:"turf_id"`
	Turf      Turf      `gorm:"foreignKey=TurfId" json:"-" binding:"-"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

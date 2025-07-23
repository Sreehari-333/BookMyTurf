package models

import (
	"time"

	"gorm.io/gorm"
)

type BlockedSlot struct {
	gorm.Model
	TurfID    uint      `json:"turf_id"`
	Turf      Turf      `json:"turf" gorm:"foreignKey:TurfID;references:ID"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Date      time.Time `json:"date"`
	Reason    string    `json:"reason"`
}

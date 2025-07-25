package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserId    uint      `json:"user_id" binding:"required"`
	User      User      `gorm:"foreignKey:UserId" json:"user" binding:"-"`
	TurfId    uint      `json:"turf_id" binding:"required"`
	Turf      Turf      `gorm:"foreignKey:TurfId" json:"turf" binding:"-"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Date      string    `json:"date" binding:"required"`
	Status    string    `json:"status" gorm:"type:enum('pending','booked','cancelled','completed');default:'pending'"`
	Payment   Payment   `gorm:"foreignKey:BookingID" json:"payment"`
}

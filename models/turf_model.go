package models

import "gorm.io/gorm"

// Creating struct for Turf

type Turf struct {
	gorm.Model
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	Price     float64 `json:"price"`
	Latitude  float64 `json:"latitude" gorm:"not null"`
	Longitude float64 `json:"longitude" gorm:"not null"`
	Address   string  `json:"address,omitempty"`
	Phone     string  `json:"phone" binding:"required" gorm:"type:varchar(15)"`
}

func (Turf) TableName() string {
	return "turfs"
}

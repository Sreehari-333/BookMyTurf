package models

import "gorm.io/gorm"

// Creating struct for Turf

type Turf struct {
	gorm.Model
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	Price     float64 `json:"price"`
	Images    string  `json:"image"`
	Latitude  float64 `json:"latitude" gorm:"not null"`
	Longitude float64 `json:"longitude" gorm:"not null"`
}

func (Turf) TableName() string {
	return "turfs"
}

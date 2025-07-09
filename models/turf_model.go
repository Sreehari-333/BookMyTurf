package models

import "gorm.io/gorm"

// Creating struct for Turf

type Turf struct {
	gorm.Model
	Name     string  `json:"name"`
	Location string  `json:"location"`
	Price    float64 `json:"price"`
	Images   string  `json:"image"`
}

func (Turf) TableName() string {
	return "turfs"
}

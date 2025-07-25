package models

import "gorm.io/gorm"

// Creating user struct

type User struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" gorm:"unique" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     int    `json:"role" gorm:"default:0"`
	Blocked  bool   `gorm:"default:false"`
	Phone    string `json:"phone" binding:"required" gorm:"type:varchar(15)"`
	City     string `json:"city" binding:"required"`
}

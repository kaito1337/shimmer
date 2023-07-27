package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login    string `gorm:"uniqueIndex"`
	Password string `gorm:"not null"`
	Name     string
	Age      uint `gorm:"check:age >= 18"`
	IsActive bool `gorm:"default:true"`
}

type UserResponse struct {
	ID       uint `json:"-"`
	Login    string
	Name     string
	Age      uint
	IsActive bool
	Password string `json:"-"`
}

var DB *gorm.DB

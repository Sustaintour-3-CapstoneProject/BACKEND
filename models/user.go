package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `gorm:"unique" json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `gorm:"unique" json:"email"`
	City        string `json:"city"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	Category    string `json:"category"`
	Picture     string `json:"picture"`
	File        string `json:"file"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
}

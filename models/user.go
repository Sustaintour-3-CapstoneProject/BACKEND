package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `gorm:"unique" json:"username"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `gorm:"unique" json:"phone_number"`
	Password    string `json:"password"`
}

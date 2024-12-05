package models

import "gorm.io/gorm"

type VideoContentView struct {
	gorm.Model
	VideoContentID uint `json:"video_content_id"`
	UserID         uint `json:"user_id"`
}

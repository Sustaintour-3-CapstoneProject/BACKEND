package models

type City struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

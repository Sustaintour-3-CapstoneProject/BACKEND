package models

type City struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

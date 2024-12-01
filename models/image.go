package models

type Image struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	DestinationID uint   `json:"destination_id"`
	URL           string `json:"url"`
}

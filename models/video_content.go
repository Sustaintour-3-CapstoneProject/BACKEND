package models

type VideoContent struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	DestinationID uint   `json:"destination_id"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	Description   string `json:"description"`
}

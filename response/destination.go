package response

import "time"

type DestinationResponse struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `json:"name"`
	City             City           `json:"city" gorm:"foreignKey:CityID;references:ID"`
	Position         float64        `json:"position"`
	Address          string         `json:"address"`
	OperationalHours string         `json:"operational_hours"`
	TicketPrice      float64        `json:"ticket_price"`
	Category         string         `json:"category"`
	Description      string         `json:"description"`
	Facilities       []string       `json:"facilities"`
	CreatedAt        time.Time      `json:"created_at"`
	Images           []Image        `json:"images" gorm:"foreignKey:DestinationID"`
	VideoContents    []VideoContent `json:"video_contents" gorm:"foreignKey:DestinationID"`
}

type City struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type Image struct {
	DestinationID uint   `json:"destination_id"`
	URL           string `json:"url"`
}

type VideoContent struct {
	DestinationID uint   `json:"destination_id"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	Description   string `json:"description"`
}

package response

import (
	"backend/models"
	"time"
)

type RouteResponse struct {
	ID                  uint                 `gorm:"primaryKey" json:"id"`
	UserID              uint                 `json:"userID"`
	OriginCityName      string               `json:"originCityName"`
	DestinationCityName string               `json:"destinationCityName"`
	Distance            float64              `json:"distance"`
	Time                string               `json:"time"`
	Cost                int                  `json:"cost"`
	CreatedAt           time.Time            `json:"created_at"`
	Destinations        []models.Destination `json:"destinations"`
}

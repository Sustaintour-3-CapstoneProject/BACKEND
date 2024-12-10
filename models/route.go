package models

import (
	"time"
)

type Route struct {
	ID                  uint               `gorm:"primaryKey" json:"id"`
	UserID              uint               `json:"userID"`
	OriginCityName      string             `json:"originCityName"`
	DestinationCityName string             `json:"destinationCityName"`
	Distance            float64            `json:"distance"`
	CreatedAt           time.Time          `json:"created_at"`
	Destinations        []RouteDestination `json:"destinations" gorm:"foreignKey:RouteID"`
}

package models

import (
	"time"
)

type RouteDestination struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	RouteID       uint      `json:"routeID"`
	DestinationID uint      `json:"destinationID"`
	CreatedAt     time.Time `json:"created_at"`
}

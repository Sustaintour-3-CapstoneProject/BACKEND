package models

import (
	"time"

	"gorm.io/gorm"
)

type Destination struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Name             string         `json:"name"`
	Address          string         `json:"address"`
	OperationalHours string         `json:"operational_hours"`
	TicketPrice      float64        `json:"ticket_price"`
	Category         string         `json:"category"`
	Description      string         `json:"description"`
	Facilities       string         `json:"facilities"`
	CreatedAt        time.Time      `json:"created_at"`
	Images           []Image        `json:"images"`   
}

func (b *Destination) AfterCreate(tx *gorm.DB) (err error) {
    return tx.Model(b).Preload("Images").Error
}
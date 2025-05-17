package entity

import (
	"github.com/google/uuid"
	"time"
)

type ProductWithQuantity struct {
	Candles  Candles `json:"candles"`
	Quantity int     `json:"quantity"`
}
type Candles struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     uint      `json:"version"`
	Price       int       `json:"price"`
	Images      []string  `json:"images"`
	ColorID     uuid.UUID `json:"color_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
	Weight      int       `json:"weight"`
}

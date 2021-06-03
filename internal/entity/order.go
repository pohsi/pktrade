package entity

import (
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	OwnerId   int       `json:"owner_id"`
	OwnerName string    `json:"owner_name"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     float64   `json:"price"`
}

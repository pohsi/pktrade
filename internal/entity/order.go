package entity

import (
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Owner     string    `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     int       `json:"price"`
}

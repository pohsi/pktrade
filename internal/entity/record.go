package entity

import (
	"time"
)

type Record struct {
	ID        int       `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     int       `json:"price"`
}

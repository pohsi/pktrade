package entity

import (
	"time"
)

type Record struct {
	ID        int       `json:"id"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     float64   `json:"price"`
}

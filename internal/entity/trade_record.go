package entity

import (
	"time"
)

type Wanted struct {
	ID        int       `json:"id"`
	Owner     int       `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     int       `json:"price"`
}

type Selling struct {
	ID        int       `json:"id"`
	User      int       `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     int       `json:"price"`
}

type Trade struct {
	ID         string    `json:"id"`
	UserFrom   int       `json:"user_from"`
	UserTo     int       `json:"user_to"`
	CompleteAt time.Time `json:"complete_at"`
}

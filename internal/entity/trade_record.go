package entity

import (
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Owner     int       `json:"owner"`
	CreatedAt time.Time `json:"created_at"`
	CardType  int       `json:"card_type"`
	Price     int       `json:"price"`
}

type PurchaseOrder struct {
	Order `json:"order"`
}

type SellOrder struct {
	Order `json:"order"`
}

type Records struct {
	PurchaseOrders []PurchaseOrder `json:"purchase_order"`
	SellOrders     []PurchaseOrder `json:"sell_order"`
}

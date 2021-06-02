package trade

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pohsi/pktrade/internal/entity"
	"github.com/pohsi/pktrade/pkg/log"
)

type Service interface {
	GetPurchaseOrder(ctx context.Context) ([]entity.PurchaseOrder, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

type CardIdType int
type OrderType int
type PriceType float64
type CardType int

const (
	Pikachu CardType = iota + 1
	Bulbasaur
	Charmander
	Squirtle
	CardTypeCount
)

const (
	OrderPurchase OrderType = iota + 1
	OrderSell
)

const (
	priceFloor    PriceType = 1
	priceCeilling PriceType = 10
)

// CreateOrderRequest reprsents purchase card by type with given price,
// we are not going to consider quantity
type CreateOrderRequest struct {
	Id        int       `json:"id"`
	OrderType OrderType `json:"order_type"`
	CardType  CardType  `json:"card_type"`
	Price     PriceType `json:"price"`
}

func (c CreateOrderRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.OrderType, validation.Required, validation.By(func(v interface{}) error {
			if val, ok := v.(OrderType); !ok || !(val == OrderPurchase || val == OrderSell) {
				return fmt.Errorf("unexcpet card id: %v", v)
			}
			return nil
		})),
	); err != nil {
		return err
	}

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.CardType, validation.Required, validation.By(func(v interface{}) error {
			if val, ok := v.(CardType); !ok || val >= CardTypeCount {
				return fmt.Errorf("unexcpet card id: %v", v)
			}
			return nil
		})),
	); err != nil {
		return err
	}

	if err := validation.ValidateStruct(&c,
		validation.Field(&c.Price, validation.Required, validation.By(func(v interface{}) error {
			if val, ok := v.(PriceType); !ok || val < priceFloor || val > priceCeilling {
				return fmt.Errorf("unexcpet price: %v", v)
			}
			return nil
		})),
	); err != nil {
		return err
	}

	return nil
}

func NewService(repo Repository, logger log.Logger) Service {
	return &service{repo, logger}
}

// Query returns recent 50 trade records for all cards.
func (s *service) GetPurchaseOrder(ctx context.Context) ([]entity.PurchaseOrder, error) {
	return nil, nil
}

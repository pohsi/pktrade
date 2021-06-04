package trade

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pohsi/pktrade/internal/entity"
	"github.com/pohsi/pktrade/pkg/log"
)

type Order struct {
	entity.Order
}

type Record struct {
	entity.Record
}

type Service interface {
	GetRecordsByCardType(c context.Context, req GetRecordsByCardTypeRequest) ([]Record, error)

	GetRecords(c context.Context, req GetRecordsRequest) (interface{}, error)

	CreateOrder(ctx context.Context, req CreateOrderRequest) (Order, error)
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

const (
	QueryPurchaseStatus = iota + 1
	QuerySellStatus
	QueryRecordStatus
)

// CreateOrderRequest reprsents purchase card by type with given price,
// we are not going to consider quantity
type CreateOrderRequest struct {
	UserId    int       `json:"user_id"`
	UserName  string    `json:"user_name"`
	OrderType OrderType `json:"order_type"`
	CardType  CardType  `json:"card_type"`
	Price     PriceType `json:"price"`
}

type GetRecordsRequest struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Limit     int    `json:"limit"`
	QueryType int    `json:"query_type"`
}

type GetRecordsByCardTypeRequest struct {
	CardType CardType `json:"card_type"`
	Limit    int      `json:"limit"`
}

func userIdValidateHelper(v interface{}) error {
	if val, ok := v.(int); !ok || !(val > 1 && val <= 100000) {
		return fmt.Errorf("unexpect user id: %v", v)
	}
	return nil
}

func cardTypeValidateHelper(v interface{}) error {
	if val, ok := v.(CardType); !ok || val <= 0 || val >= CardTypeCount {
		return fmt.Errorf("unexpect card type: %v", v)
	}
	return nil
}

func (c CreateOrderRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.UserId, validation.Required, validation.By(userIdValidateHelper)),
		validation.Field(&c.UserName, validation.Required, validation.Length(1, 128)),
		validation.Field(&c.OrderType, validation.Required, validation.By(func(v interface{}) error {
			if val, ok := v.(OrderType); !ok || !(val == OrderPurchase || val == OrderSell) {
				return fmt.Errorf("unexpect order type: %v", v)
			}
			return nil
		})),
		validation.Field(&c.CardType, validation.Required, validation.By(cardTypeValidateHelper)),
		validation.Field(&c.Price, validation.Required, validation.By(func(v interface{}) error {
			if val, ok := v.(PriceType); !ok || val < priceFloor || val > priceCeilling {
				return fmt.Errorf("unexpect price: %v", v)
			}
			return nil
		})),
	); err != nil {
		return err
	}

	return nil
}

func (c GetRecordsRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.UserId, validation.Required, validation.By(userIdValidateHelper)),
		validation.Field(&c.UserName, validation.Required, validation.Length(1, 128)),
		validation.Field(&c.QueryType, validation.Required, validation.In(QueryPurchaseStatus,
			QuerySellStatus,
			QueryRecordStatus)),
	); err != nil {
		return err
	}

	return nil
}

func (c GetRecordsByCardTypeRequest) Validate() error {
	if err := validation.ValidateStruct(&c,
		validation.Field(&c.CardType, validation.Required, validation.By(cardTypeValidateHelper)),
	); err != nil {
		return err
	}

	return nil
}

func NewService(repo Repository, logger log.Logger) Service {
	return &service{repo, logger}
}

func (s *service) getUserRecords(c context.Context, req GetRecordsRequest) ([]Record, error) {
	r, err := s.repo.GetUserRecords(c, req.UserName, req.Limit)
	if err != nil {
		return nil, err
	}
	result := []Record{}
	for _, item := range r {
		result = append(result, Record{item})
	}
	return result, nil
}

func (s *service) getUserPurchase(c context.Context, req GetRecordsRequest) ([]Order, error) {
	r, err := s.repo.GetUserPurchaseOrders(c, req.UserName, req.Limit)
	if err != nil {
		return nil, err
	}
	result := []Order{}
	for _, item := range r {
		result = append(result, Order{item})
	}
	return result, nil
}

func (s *service) getUserSell(c context.Context, req GetRecordsRequest) ([]Order, error) {
	r, err := s.repo.GetUserSellOrders(c, req.UserName, req.Limit)
	if err != nil {
		return nil, err
	}
	result := []Order{}
	for _, item := range r {
		result = append(result, Order{item})
	}
	return result, nil
}

// Query returns recent 50 trade records for all cards.
func (s *service) GetRecords(c context.Context, req GetRecordsRequest) (interface{}, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	switch req.QueryType {
	case QueryRecordStatus:
		return s.getUserRecords(c, req)
	case QuerySellStatus:
		return s.getUserSell(c, req)
	case QueryPurchaseStatus:
		return s.getUserPurchase(c, req)
	}
	return nil, fmt.Errorf("Unknow query type: %v", req.QueryType)
}

func (s *service) GetRecordsByCardType(c context.Context, req GetRecordsByCardTypeRequest) ([]Record, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	r, err := s.repo.GetRecordsByCardType(c, int(req.CardType), req.Limit)
	if err != nil {
		return nil, err
	}
	result := []Record{}
	for _, item := range r {
		result = append(result, Record{item})
	}
	return result, nil
}

func (s *service) CreateOrder(ctx context.Context, req CreateOrderRequest) (Order, error) {
	if err := req.Validate(); err != nil {
		return Order{}, err
	}

	order := entity.Order{
		OwnerId:   req.UserId,
		OwnerName: req.UserName,
		CreatedAt: time.Now(),
		CardType:  int(req.CardType),
		Price:     float64(req.Price),
	}

	if req.OrderType == OrderSell {
		if err := s.repo.ResolverOrderSell(ctx, order); err != nil {
			return Order{}, err
		}
	} else if err := s.repo.ResolverOrderPurchase(ctx, order); err != nil {
		return Order{}, err
	}

	return Order{order}, nil
}

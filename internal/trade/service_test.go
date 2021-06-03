package trade

import (
	"context"
	"errors"
	"testing"

	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"
)

var errorCRUD = errors.New("error crud")

func TestCreateOrderRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateOrderRequest
		wantError bool
	}{
		{name: "Success", model: CreateOrderRequest{UserId: 33, UserName: "user33", OrderType: OrderPurchase, CardType: Pikachu, Price: 9.9}, wantError: false},
		{name: "Success", model: CreateOrderRequest{UserId: 168, UserName: "user168", OrderType: OrderSell, CardType: Bulbasaur, Price: 7.45}, wantError: false},
		{name: "BadOrderType", model: CreateOrderRequest{UserId: 964, UserName: "user964", OrderType: 3, CardType: Squirtle, Price: 4.85}, wantError: true},
		{name: "BadCardId", model: CreateOrderRequest{UserId: 12, UserName: "user12", OrderType: OrderPurchase, CardType: 0, Price: 3.3}, wantError: true},
		{name: "BadPriceToLow", model: CreateOrderRequest{UserId: 57, UserName: "user57", OrderType: OrderSell, CardType: Charmander, Price: 0.8}, wantError: true},
		{name: "BadPriceToHigh", model: CreateOrderRequest{UserId: 88, UserName: "user88", OrderType: OrderSell, CardType: CardTypeCount, Price: 10.5}, wantError: true},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			err := v.model.Validate()
			t.Logf("%v => Error: %v", v, err)
			assert.Equal(t, v.wantError, err != nil)
		})
	}
}

func Test_serviceCRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRespository{}, logger)

	ctx := context.Background()

	s.GetPurchaseOrder(ctx)

}

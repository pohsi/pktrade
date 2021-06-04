package trade

import (
	"context"
	"testing"

	"github.com/pohsi/pktrade/internal/test"
	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"

	"github.com/pohsi/pktrade/internal/entity"
)

type mockRespository struct {
	purchaseOrder []entity.Order
	sellOrder     []entity.Order
	records       []entity.Record
}

func (m *mockRespository) GetUserRecords(c context.Context, userName string, limit int) ([]entity.Record, error) {
	return m.records, nil
}

func (m *mockRespository) GetUserSellOrders(c context.Context, userName string, limit int) ([]entity.Order, error) {
	return nil, nil
}

func (m *mockRespository) GetUserPurchaseOrders(c context.Context, userName string, limit int) ([]entity.Order, error) {
	return nil, nil
}

func (m *mockRespository) ResolverOrderSell(c context.Context, order entity.Order) error {
	return nil
}

func (m *mockRespository) GetRecordsByCardType(c context.Context, cardType int, limit int) ([]entity.Record, error) {
	return nil, nil
}

func (m *mockRespository) ResolverOrderPurchase(c context.Context, order entity.Order) error {
	return nil
}

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	assert.NotNil(t, logger)

	test.MockRouter(logger)

}

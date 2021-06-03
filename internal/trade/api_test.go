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

func (m *mockRespository) GetPurchaseOrder() ([]Order, error) {
	return m.GetPurchaseOrder()
}

func (m *mockRespository) ResolverOrderSell(c context.Context, order entity.Order) error {
	return nil
}

func (m *mockRespository) ResolverOrderPurchase(c context.Context, order entity.Order) error {
	return nil
}

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	assert.NotNil(t, logger)

	test.MockRouter(logger)

}

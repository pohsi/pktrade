package trade

import (
	"testing"

	"github.com/pohsi/pktrade/internal/test"
	"github.com/pohsi/pktrade/pkg/log"
	"github.com/stretchr/testify/assert"

	"github.com/pohsi/pktrade/internal/entity"
)

type mockRespository struct {
	purchaseOrder []entity.PurchaseOrder
	sellOrder     []entity.SellOrder
	records       []entity.Records
}

func (m *mockRespository) GetPurchaseOrder() ([]entity.PurchaseOrder, error) {
	return m.GetPurchaseOrder()
}

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	assert.NotNil(t, logger)

	test.MockRouter(logger)

}

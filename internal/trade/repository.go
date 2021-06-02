package trade

import (
	"github.com/pohsi/pktrade/internal/entity"
	"github.com/pohsi/pktrade/pkg/dbconnection"
	"github.com/pohsi/pktrade/pkg/log"
)

type Repository interface {
	GetPurchaseOrder() ([]entity.PurchaseOrder, error)
}

type repository struct {
	db     dbconnection.DB
	logger log.Logger
}

func NewRepository(db dbconnection.DB, logger log.Logger) Repository {
	return &repository{db, logger}
}

func (r *repository) GetPurchaseOrder() ([]entity.PurchaseOrder, error) {
	// var count int
	// err := r.db.With(ctx).Select("COUNT(*)").From("album").Row(&count)
	return []entity.PurchaseOrder{}, nil
}

package trade

import (
	"github.com/pohsi/pktrade/pkg/dbconnection"
	"github.com/pohsi/pktrade/pkg/log"
)

type Repository interface {
	GetPurchaseOrder() ([]Order, error)
}

type repository struct {
	db     dbconnection.DB
	logger log.Logger
}

func NewRepository(db dbconnection.DB, logger log.Logger) Repository {
	return &repository{db, logger}
}

func (r *repository) GetPurchaseOrder() ([]Order, error) {
	// var count int
	// err := r.db.With(ctx).Select("COUNT(*)").From("album").Row(&count)
	return []Order{}, nil
}

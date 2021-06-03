package trade

import (
	"context"
	"database/sql"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/pohsi/pktrade/internal/entity"
	"github.com/pohsi/pktrade/pkg/dbconnection"
	"github.com/pohsi/pktrade/pkg/log"
)

type Repository interface {
	GetPurchaseOrder() ([]Order, error)

	ResolverOrderSell(c context.Context, order entity.Order) error

	ResolverOrderPurchase(c context.Context, order entity.Order) error
}

const (
	cardOrderSellTable     string = "card_sell"
	cardOrderPurchaseTable string = "card_purchase"
)

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

func (r *repository) ResolverOrderPurchase(c context.Context, order entity.Order) error {
	err := r.db.Transactional(c, func(c context.Context) error {
		tx := r.db.With(c)
		var wanted entity.Order
		err := tx.Select().
			From(cardOrderSellTable).
			Where(dbx.And(
				dbx.NewExp(`price>={:price}`, dbx.Params{"price": order.Price}),
				dbx.NewExp(`card_type={:cardType}`, dbx.Params{"cardType": order.CardType}),
			)).
			OrderBy("price ASC").
			One(&wanted)

		if err != nil && err != sql.ErrNoRows {
			return err
		}
		r.logger.Info(wanted)
		return nil
	})

	return err
}

func (r *repository) ResolverOrderSell(c context.Context, order entity.Order) error {
	return nil

}

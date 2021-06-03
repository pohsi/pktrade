package trade

import (
	"context"
	"database/sql"
	"time"

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

type card_purchase struct {
	entity.Order
}

type card_sell struct {
	entity.Order
}

func (r *repository) ResolverOrderPurchase(c context.Context, order entity.Order) error {
	err := r.db.Transactional(c, func(c context.Context) error {
		tx := r.db.With(c)
		var wanted entity.Order
		err := tx.Select().
			From(cardOrderSellTable).
			Where(dbx.And(
				// Should we allow purchase from ourself by check OwnerId?
				dbx.NewExp(`price<={:price}`, dbx.Params{"price": order.Price}),
				dbx.NewExp(`card_type={:cardType}`, dbx.Params{"cardType": order.CardType}),
			)).
			OrderBy("price ASC").
			One(&wanted)

		if err != nil {
			return err
		}

		tx.Model(&card_sell{wanted}).Delete()
		record := entity.Record{
			FromUser:  wanted.OwnerName,
			ToUser:    order.OwnerName,
			CreatedAt: time.Now(),
			CardType:  order.CardType,
			Price:     wanted.Price,
		}
		tx.Model(&record).Insert()
		return nil
	})

	if err == sql.ErrNoRows {
		err = r.postPurchaseOrder(c, card_purchase{order})
		return err
	}

	return err
}

func (r *repository) postPurchaseOrder(c context.Context, order card_purchase) error {
	return r.db.With(c).Model(&order).Insert()
}

func (r *repository) postSellOrder(c context.Context, order card_sell) error {
	return r.db.With(c).Model(&order).Insert()
}

func (r *repository) ResolverOrderSell(c context.Context, order entity.Order) error {
	err := r.db.Transactional(c, func(c context.Context) error {
		tx := r.db.With(c)
		var wanted entity.Order
		err := tx.Select().
			From(cardOrderPurchaseTable).
			Where(dbx.And(
				// Should we allow sell to ourself?
				dbx.NewExp(`price>={:price}`, dbx.Params{"price": order.Price}),
				dbx.NewExp(`card_type={:cardType}`, dbx.Params{"cardType": order.CardType}),
			)).
			OrderBy("price DESC").
			One(&wanted)

		if err != nil {
			return err
		}

		if err = tx.Model(&card_purchase{wanted}).Delete(); err != nil {
			return err
		}
		record := entity.Record{
			FromUser:  order.OwnerName,
			ToUser:    wanted.OwnerName,
			CreatedAt: time.Now(),
			CardType:  order.CardType,
			Price:     wanted.Price,
		}
		if err = tx.Model(&record).Insert(); err != nil {
			return err
		}
		return nil
	})

	if err == sql.ErrNoRows {
		err = r.postSellOrder(c, card_sell{order})
		return err
	}

	return err
}

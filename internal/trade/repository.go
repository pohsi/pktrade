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
	GetUserRecords(c context.Context, userName string, limit int) ([]entity.Record, error)

	GetUserSellOrders(c context.Context, userName string, limit int) ([]entity.Order, error)

	GetUserPurchaseOrders(c context.Context, userName string, limit int) ([]entity.Order, error)

	GetRecordsByCardType(c context.Context, cardType int, limit int) ([]entity.Record, error)

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

// GetRecords queries order records where user name equal to from_user or to_user
func (r *repository) GetUserRecords(c context.Context, userName string, limit int) ([]entity.Record, error) {
	var records []entity.Record
	err := r.db.With(c).
		Select().
		OrderBy().
		From("record").
		Where(dbx.Or(
			dbx.NewExp(`from_user={:name}`, dbx.Params{"name": userName}),
			dbx.NewExp(`to_user={:name}`, dbx.Params{"name": userName}),
		)).
		OrderBy("created_at DESC").
		Limit(int64(limit)).
		All(&records)

	return records, err
}

func (r *repository) GetUserSellOrders(c context.Context, userName string, limit int) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.With(c).
		Select().
		OrderBy().
		From("card_sell").
		Where(dbx.Or(
			dbx.NewExp(`owner_name={:name}`, dbx.Params{"name": userName}),
		)).
		OrderBy("created_at DESC").
		Limit(int64(limit)).
		All(&orders)

	return orders, err
}

func (r *repository) GetUserPurchaseOrders(c context.Context, userName string, limit int) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.With(c).
		Select().
		OrderBy().
		From("card_purchase").
		Where(dbx.Or(
			dbx.NewExp(`owner_name={:name}`, dbx.Params{"name": userName}),
		)).
		OrderBy("created_at DESC").
		Limit(int64(limit)).
		All(&orders)

	return orders, err
}

func (r *repository) GetRecordsByCardType(c context.Context, cardType int, limit int) ([]entity.Record, error) {
	var records []entity.Record
	err := r.db.With(c).
		Select().
		OrderBy().
		From("record").
		Where(
			dbx.NewExp(`card_type={:type}`, dbx.Params{"type": cardType}),
		).
		OrderBy("created_at DESC").
		Limit(int64(limit)).
		All(&records)

	return records, err
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

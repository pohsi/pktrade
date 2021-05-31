package dbconnection

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

type DB interface {
	Transactional(ctx context.Context, f func(ctx context.Context) error) error

	With(ctx context.Context) dbx.Builder
}

type contextKey int

const (
	txKey contextKey = iota
)

type concreteDB struct {
	db *dbx.DB
}

func New(db *dbx.DB) DB {
	return &concreteDB{db}
}

func (db *concreteDB) DB() *dbx.DB {
	return db.db
}

func (db *concreteDB) With(ctx context.Context) dbx.Builder {
	if tx, ok := ctx.Value(txKey).(*dbx.Tx); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

func (db *concreteDB) Transactional(ctx context.Context, f func(ctx context.Context) error) error {
	return db.db.TransactionalContext(ctx, nil, func(tx *dbx.Tx) error {
		return f(context.WithValue(ctx, txKey, tx))
	})
}

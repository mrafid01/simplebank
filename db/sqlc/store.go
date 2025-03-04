package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, args CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, args VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// Store provide all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRoll := tx.Rollback(); errRoll != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, errRoll)
		}
		return err
	}

	return tx.Commit()
}

package db

import (
	"context"
	"fmt"
)

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if errRoll := tx.Rollback(ctx); errRoll != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, errRoll)
		}
		return err
	}

	return tx.Commit(ctx)
}

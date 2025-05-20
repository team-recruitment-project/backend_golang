package repository

import (
	"backend_golang/ent"
	"context"
)

// TransactionManager handles database transactions
type TransactionManager struct {
	client *ent.Client
}

// NewTransactionManager creates a new TransactionManager
func NewTransactionManager(client *ent.Client) *TransactionManager {
	return &TransactionManager{
		client: client,
	}
}

// WithTx executes the given function within a transaction
func (tm *TransactionManager) WithTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := tm.client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = rerr
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

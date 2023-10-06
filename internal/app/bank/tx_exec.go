package bank

import (
	"context"
	"fmt"

	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/db"
)

// execTx executes a callback function within a database transaction, finally it commits or rollbacks the transaction.
func (bank *SQLBank) execTx(ctx context.Context, fn func(*db.Queries) error) error {
	// Begin transaction
	tx, err := bank.connPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("exec tx:create transaction: %v", err)
	}

	q := db.New(tx)

	// Execute transaction
	if err := fn(q); err != nil {
		// Rollback transaction
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("exec tx:transaction err: %v, rollback err: %v", err, rbErr)
		}
		return fmt.Errorf("exec tx:transaction err: %v", err)
	}

	// Commit transaction
	return tx.Commit(ctx)
}

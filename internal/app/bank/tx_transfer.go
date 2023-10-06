package bank

import (
	"context"

	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/db"
)

// TransferParams contains the input parameters of the transfer transaction
type TransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferResult is the result of the transfer transaction
type TransferResult struct {
	// Created transfer record
	Transfer db.Transfer `json:"transfer"`
	// From account after its balance is updated
	FromAccount db.Account `json:"from_account"`
	// To account after its balance is updated
	ToAccount db.Account `json:"to_account"`
	// From entry records that money is moving out
	FromEntry db.Entry `json:"from_entry"`
	// To entry records that money is moving in
	ToEntry db.Entry `json:"to_entry"`
}

// Transfer performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction.
func (bank *SQLBank) Transfer(ctx context.Context, transfer TransferParams) (TransferResult, error) {
	var result TransferResult

	// Create a transaction with the callback function
	err := bank.execTx(ctx, func(q *db.Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, db.CreateTransferParams{
			FromAccountID: transfer.FromAccountID,
			ToAccountID:   transfer.ToAccountID,
			Amount:        transfer.Amount, // Amount to transfer
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: transfer.FromAccountID,
			Amount:    -transfer.Amount, // Money moves out from account
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, db.CreateEntryParams{
			AccountID: transfer.ToAccountID,
			Amount:    transfer.Amount, // Money moves in to account
		})
		if err != nil {
			return err
		}

		// Some notes about database locks.
		// Its always good to be consistent in the way database locks should be handled.
		// To handle dead locks we make sure to apply database locks in a consistent order, in
		// our case we always update accounts with smaller ids first.
		if transfer.FromAccountID < transfer.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx,
				q,
				transfer.FromAccountID,
				-transfer.Amount,
				transfer.ToAccountID,
				transfer.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, transfer.ToAccountID, transfer.Amount, transfer.FromAccountID, -transfer.Amount)
		}

		return err
	})

	return result, err
}

// addMoney add/withdraw money to the two accounts
func addMoney(
	ctx context.Context,
	q *db.Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 db.Account, account2 db.Account, err error) {
	account1, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}

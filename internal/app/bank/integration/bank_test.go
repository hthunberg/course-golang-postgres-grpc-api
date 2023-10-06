//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/bank"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/db"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/pkg/currency"
	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/pkg/random"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, user db.User, currency string) db.Account {
	ctx := context.Background()

	accParams := db.CreateAccountParams{
		Owner:    user.Username,
		Balance:  random.Money(),
		Currency: currency,
	}

	account, err := testee.CreateAccount(ctx, accParams)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, accParams.Owner, account.Owner)
	require.Equal(t, accParams.Balance, account.Balance)
	require.Equal(t, accParams.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createRandomUser(t *testing.T) db.User {
	ctx := context.Background()

	userParams := db.CreateUserParams{
		Username:       random.Owner(),
		HashedPassword: random.String(10),
		FullName:       random.Owner(),
		Email:          random.Email(),
	}

	user, err := testee.CreateUser(ctx, userParams)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, userParams.Email, user.Email)
	require.Equal(t, userParams.FullName, user.FullName)
	require.Equal(t, userParams.HashedPassword, user.HashedPassword)
	require.Equal(t, userParams.Username, user.Username)

	require.NotEmpty(t, user.Username)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestAccountTransfer(t *testing.T) {
	// TODO: We might need to truncate between test

	// Accounts need to have the same currency, we do not support exchange
	// rates in this demo bank
	account1 := createRandomAccount(t, createRandomUser(t), currency.SEK)
	account2 := createRandomAccount(t, createRandomUser(t), currency.SEK)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	// Channels to get results from the goroutines running transactions
	errs := make(chan error)
	results := make(chan bank.TransferResult)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := testee.Transfer(context.Background(), bank.TransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			// Communicate errors and result to the surrounding goroutine
			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testee.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testee.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testee.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedAccount1, err := testee.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testee.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)
}

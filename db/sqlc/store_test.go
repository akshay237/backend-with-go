package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	// 0. create store object
	store := NewStore(testDB)

	// 1. create two accounts for transfer
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Before transaction: ", account1.Balance, account2.Balance)

	// 2. create n concurrent transfers b/w accounts and channel of err and transfer txn
	n := 5
	amount := int64(10)

	txnResults := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			txnResults <- result

		}()
	}

	// 3. loop over the concurrent operations and check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-txnResults
		require.NotEmpty(t, result)

		// check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		// from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// ToDO: check for account balance updates
		// check from account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check accounts balance
		fmt.Println("After transaction: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
	}

	// check the final updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updateAccount1.Balance, account1.Balance-int64(n)*amount)
	require.Equal(t, updatedAccount2.Balance, account2.Balance+int64(n)*amount)
}

func TestTransferTxDeadLock(t *testing.T) {
	// 0. create store object
	store := NewStore(testDB)

	// 1. create two accounts for transfer
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("Before transaction: ", account1.Balance, account2.Balance)

	// 2. create n concurrent transfers b/w accounts and channel of err and transfer txn
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})

			errs <- err
		}()
	}

	// 3. loop over the concurrent operations and check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	// check the final updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, updateAccount1.Balance, account1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)
	fmt.Println("After transaction: ", updateAccount1.Balance, updatedAccount2.Balance)
}

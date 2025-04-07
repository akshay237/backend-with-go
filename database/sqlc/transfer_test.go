package database

import (
	"context"
	"testing"

	"github.com/akshay237/backend-with-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {

	// 1. create the args for transfer
	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	// 2. calls the create transfer db function
	transfer, err := testQueries.CreateTransfer(context.Background(), args)

	// 3. check the err and other properties are not nil
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, transfer.ToAccountID, args.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {

	// 1. create two accounts used for transfer
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// calls create random Transfer
	createRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {

	// 1. create the accounts first
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// 2. create a random transfer
	transfer1 := createRandomTransfer(t, account1, account2)

	// 3. get the transfer
	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	// 4. check for no error and other property
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
}

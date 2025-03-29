package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/akshay237/backend-with-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	// 1. create the args passed to the query
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	// 2. calls the create Account function generated through sqlc
	account, err := testQueries.CreateAccount(context.Background(), args)

	// 3. Check for no error and other properties using testify library
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {

	// 1. create an account
	account1 := createRandomAccount(t)

	// 2. pass the account id of account got from create account to get Account
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	// 3. check the err is nil and other properties are same
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {

	// 1. Create an account
	account1 := createRandomAccount(t)

	// 2. Update the balance of the created account
	args := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), args)

	// 3. check error for not nil and other properties
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, args.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
}

func TestDeleteAccount(t *testing.T) {

	// 1. create an account
	account1 := createRandomAccount(t)

	// 2. delete the created account
	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	// 3. check for error is not nil
	require.NoError(t, err)

	// 4. get the account and it will return an error
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {

	// 1. create 5 random accounts
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}

	// 2. create the args for list accounts query
	args := ListAccountsParams{
		Limit:  3,
		Offset: 2,
	}

	// 3. list the accounts
	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 3)

	// 4. loop over the accounts and check each one is not empty
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

package database

import (
	"context"
	"testing"
	"time"

	"github.com/akshay237/backend-with-go/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) Entry {

	// 1. create the params for Create Entry query
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	// 2. calls the create entry db function
	entry, err := testQueries.CreateEntry(context.Background(), args)

	// 3. Check for error is nil
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.Equal(t, args.Amount, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	// 1. create an account first
	account := createRandomAccount(t)

	// 2. create an entry
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {

	// 1. create an account first
	account := createRandomAccount(t)

	// 2. create an entry
	entry1 := createRandomEntry(t, account)

	// 3. get the entry for above created entry
	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	// 4. check for error nil and other properties are same
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {

	// 1. create an account first
	account := createRandomAccount(t)

	// 2. create random entries
	for i := 0; i < 5; i++ {
		createRandomEntry(t, account)
	}

	// 3. create args for list entries query
	args := ListEntriesParams{
		AccountID: account.ID,
		Limit:     3,
		Offset:    2,
	}

	// 4. calls the list entries db function
	entries, err := testQueries.ListEntries(context.Background(), args)

	// 5. check for no error
	require.NoError(t, err)
	require.Len(t, entries, 3)

	// 5. loop over the entries and check they are not empty
	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}

package database

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore returns a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	// 1. Begin the transaction
	txn, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 2. create the queries object by passing the txn instance
	q := New(txn)
	err = fn(q)
	if err != nil {
		if rbErr := txn.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// 3. if there are no errors commit the txn
	return txn.Commit()
}

// TransferTxParams to perform a transfer b/w accounts
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult to store the result of this txn
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfers from one account to the other account.
// It creates a transfer record, add account entries and update accounts balance with in a single transaction.
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	// 1. create a var of tx result
	var result TransferTxResult

	// 2. calls the execTxn
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		// 2.1 first create a transfer by calling the queries create transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2.2 create an entry in from account for balance debited
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// 2.3 create an entry in to account for balance credited
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// ToDO: update the account's balance
		return nil
	})

	return result, err
}

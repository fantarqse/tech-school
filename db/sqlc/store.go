package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db quries & transactions.
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore create a new Store.
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	if err := fn(q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error:%v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// TransferTxParams contains the input parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccountID sql.NullInt64 `json:"from_account_id"`
	ToAccountID   sql.NullInt64 `json:"to_account_id"`
	Amount        int64         `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, adds account entries, and updates accounts' balance within a single database transaction.
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	if err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID.Int64 < arg.ToAccountID.Int64 {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID.Int64, -arg.Amount, arg.ToAccountID.Int64, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID.Int64, arg.Amount, arg.FromAccountID.Int64, -arg.Amount)
		}
		if err != nil {
			return err
		}

		return nil

	}); err != nil {
		return TransferTxResult{}, err
	}

	return result, nil
}

func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return Account{}, Account{}, err
	}

	return account1, account2, nil
}

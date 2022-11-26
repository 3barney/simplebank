package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db Queries and transactions,
// It extends Queries struct in db.go via composition
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a Store instance
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// Tx takes a context and a callback function and starts a new db transaction
// create a new query object with that transaction and call callback function with created queries
// and finally commit or rollback transaction based on error returned by the transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	// Pass nil to use default db Isolation level
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// use New function and pass in db.Transaction object
	queries := New(tx)

	// Call callback fn on queries and if errored do rollback
	err = fn(queries)
	if err != nil {
		// Rollback (Rollback also raises an error: Do combine the two)
		if rollbackError := tx.Rollback(); rollbackError != nil {
			return fmt.Errorf("tx err: %v, rollback err: %v", err, rollbackError)
		}
		// if rollback is successfull: Return original Tx err
		return err
	}

	return tx.Commit()
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx transfers money from one account to another,
// 1. Creates a transfer record,
// 2. add account entries (from and to),
// 3. Update account balances accordingly within a single db transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		// NOTE: we access result and arg objects of outer function, making the callback function become a closure,
		// it is used to get result from callback function since a callback doesn't know the exact type of result it should return
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // Subtract since money is moving out of this acc
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID { // This prevents a deadlock by ensuring consistency in acquiring locks, by ordering account id's

			// if fromAccountId < toAccId: we update from account to account
			result.FromAccount, result.ToAccount, err = addMoney(ctx, queries, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, queries, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	//return with no params is ok since we are using Named return variables
	return
}

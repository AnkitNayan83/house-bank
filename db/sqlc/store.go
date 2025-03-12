package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// store provides all the functions to execute db queries and transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore create a new store from store sturct
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transeaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(context.Background()); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(context.Background())
}

type TransferMoneyTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransfeMoneyTxResult struct {
	Transfer    *Transfer `json:"transfer"`
	FromAccount *Account  `json:"from_account"`
	ToAccount   *Account  `json:"to_account"`
	FromEntry   *Entry    `json:"from_entry"`
	ToEntry     *Entry    `json:"to_entry"`
}

func (store *Store) TransferMoneyTx(ctx context.Context, arg TransferMoneyTxParams) (TransfeMoneyTxResult, error) {
	var result TransfeMoneyTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		// create transfer
		transfer, err := q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}
		result.Transfer = &transfer

		// create two entries
		// from entry
		fromEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}
		result.FromEntry = &fromEntry

		// to entry
		toEntry, err := q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}
		result.ToEntry = &toEntry

		//TODO: ADD LOCKS
		// // Get Account for from account
		// fromAccount, err := q.GetAccountById(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }

		// if fromAccount.Balance < arg.Amount {
		// 	return fmt.Errorf("insufficient balance in %v's account", fromAccount.Owner)
		// }

		// // Update Balance for from account
		// updatedFromAccount, err := q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: fromAccount.Balance - arg.Amount,
		// })

		// if err != nil {
		// 	return err
		// }
		// result.FromAccount = &updatedFromAccount

		// // Get Account for to account
		// toAccount, err := q.GetAccountById(ctx, arg.ToAccountID)
		// if err != nil {
		// 	return err
		// }

		// // Update Balance for to account
		// updatedToAccount, err := q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.ToAccountID,
		// 	Balance: toAccount.Balance + arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }
		// result.ToAccount = &updatedToAccount

		return nil
	})

	return result, err
}

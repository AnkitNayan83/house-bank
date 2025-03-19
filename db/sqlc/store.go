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

// txKey is a custom key for transaction context. It will allow us to pass the name of the transaction
// to the context so we can log it later
// type txKetType struct{}

// var txKey = txKetType{}

func (store *Store) TransferMoneyTx(ctx context.Context, arg TransferMoneyTxParams) (TransfeMoneyTxResult, error) {
	var result TransfeMoneyTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		// txName := ctx.Value(txKey)

		// fmt.Println(txName, ">> create transfer")
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

		fromAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromAccount = &fromAccount

		toAccount, err := q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: +arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToAccount = &toAccount

		return nil
	})

	return result, err
}

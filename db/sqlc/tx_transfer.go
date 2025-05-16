package db

import "context"

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

func (store *SQLStore) TransferMoneyTx(ctx context.Context, arg TransferMoneyTxParams) (TransfeMoneyTxResult, error) {
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
		// to avaoid dl
		if arg.FromAccountID < arg.ToAccountID {

			fromAcc, toAcc, err := addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

			if err != nil {
				return err
			}

			result.FromAccount = &fromAcc
			result.ToAccount = &toAcc
		} else {
			toAcc, fromAcc, err := addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)

			if err != nil {
				return err
			}

			result.FromAccount = &fromAcc
			result.ToAccount = &toAcc
		}

		return nil
	})

	return result, err
}

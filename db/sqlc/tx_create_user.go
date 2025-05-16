package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreateUser func(user User) error // callback func to send task to queue
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		user, err := q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		result.User = user

		return arg.AfterCreateUser(user)
	})

	return result, err
}

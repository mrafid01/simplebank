package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

// CreateUserTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, args CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, args.CreateUserParams)
		if err != nil {
			return err
		}

		return args.AfterCreate(result.User)
	})
	return result, err
}

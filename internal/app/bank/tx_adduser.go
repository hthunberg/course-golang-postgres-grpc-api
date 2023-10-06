package bank

import (
	"context"

	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/db"
)

type AddUserParams struct {
	db.CreateUserParams
	AfterCreate func(user db.User) error
}

type AddUserResult struct {
	User db.User `json:"user"`
}

func (store *SQLBank) AddUser(ctx context.Context, arg AddUserParams) (AddUserResult, error) {
	var result AddUserResult

	err := store.execTx(ctx, func(q *db.Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}

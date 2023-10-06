package bank

import (
	"context"

	"github.com/hthunberg/course-golang-postgres-grpc-api/internal/app/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Bank defines all functions to execute db queries and transactions
type Bank interface {
	db.Querier
	Transfer(ctx context.Context, arg TransferParams) (TransferResult, error)
	AddUser(ctx context.Context, arg AddUserParams) (AddUserResult, error)
}

// SQLBank a composition that provides transactions over multiple database queries.
// Composition allows you to build complex types by combining simpler types, promoting
// code modularity and flexibility.
type SQLBank struct {
	*db.Queries
	connPool *pgxpool.Pool
}

func NewBank(connPool *pgxpool.Pool) *SQLBank {
	return &SQLBank{
		connPool: connPool,
		Queries:  db.New(connPool),
	}
}

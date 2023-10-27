package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	// PostgreSQL Error Codes
	// https://www.postgresql.org/docs/current/errcodes-appendix.html
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

var ErrForeignKeyViolation = &pgconn.PgError{
	Code: ForeignKeyViolation,
}

// ErrorCode returns the db error code matching the first error in err's tree that matches.
func ErrorCode(err error) string {
	var pgErr *pgconn.PgError

	// Find the first error in err's tree that matches target
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

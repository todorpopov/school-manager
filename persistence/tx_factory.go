package persistence

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ITransactionFactory interface {
	BeginTransaction(ctx context.Context) (pgx.Tx, *exceptions.AppError)
	CommitOrRollback(ctx context.Context, tx pgx.Tx, err *exceptions.AppError) *exceptions.AppError
}

type TransactionFactory struct {
	db *Database
}

func NewTransactionFactory(db *Database) *TransactionFactory {
	return &TransactionFactory{db: db}
}

func (tf *TransactionFactory) BeginTransaction(ctx context.Context) (pgx.Tx, *exceptions.AppError) {
	tx, err := tf.db.Pool.Begin(ctx)
	if err != nil {
		return nil, exceptions.PgErrorToAppError(err)
	}
	return tx, nil
}

func (tf *TransactionFactory) CommitOrRollback(ctx context.Context, tx pgx.Tx, err *exceptions.AppError) *exceptions.AppError {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return err
		}
		return err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return exceptions.PgErrorToAppError(commitErr)
	}

	return nil
}

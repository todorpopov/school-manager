package domain_model

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

func PgErrorToAppError(err error) *exceptions.AppError {
	appErr := exceptions.AppError{Code: "DATABASE_ERROR", Message: "Database Error", Cause: err}

	var pgxErr *pgconn.PgError
	if !errors.As(err, &pgxErr) {
		return &appErr
	}

	switch pgxErr.Code {
	// Class 23 - Integrity constraint violation
	case "23000":
		appErr.Code = "INTEGRITY_CONSTRAINT_VIOLATION"
		appErr.Message = "Integrity constraint violation"
		return &appErr
	case "23001":
		appErr.Code = "RESTRICT_VIOLATION"
		appErr.Message = "Foreign key violation"
		return &appErr
	case "23002":
		appErr.Code = "NOT_NULL_VIOLATION"
		appErr.Message = "Not null violation"
		return &appErr
	case "23503":
		appErr.Code = "FOREIGN_KEY_VIOLATION"
		appErr.Message = "Foreign key violation"
		return &appErr
	case "23505":
		appErr.Code = "UNIQUE_VIOLATION"
		appErr.Message = "Unique violation"
		return &appErr
	default:
		return &appErr
	}
}

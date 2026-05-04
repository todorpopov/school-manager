package terms

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ITermRepository interface {
	CreateTerm(ctx context.Context, tx pgx.Tx, createTerm *CreateTerm) (*Term, *exceptions.AppError)
	GetTermById(ctx context.Context, tx pgx.Tx, termId int32) (*Term, *exceptions.AppError)
	GetTerms(ctx context.Context, tx pgx.Tx) ([]Term, *exceptions.AppError)
	DeleteTerm(ctx context.Context, tx pgx.Tx, termId int32) *exceptions.AppError
}

type TermRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewTermRepository(db *persistence.Database, logger *zap.Logger) *TermRepository {
	return &TermRepository{db, logger}
}

func (tr *TermRepository) CreateTerm(ctx context.Context, tx pgx.Tx, createTerm *CreateTerm) (*Term, *exceptions.AppError) {
	sql := `
		INSERT INTO terms (name, start_date, end_date)
		VALUES ($1, $2, $3)
		RETURNING term_id, name, start_date, end_date;
	`

	var term Term
	var err error

	if tx != nil {
		tr.logger.Debug("Creating term in transaction")
		err = tx.QueryRow(ctx, sql, createTerm.Name, createTerm.StartDate, createTerm.EndDate).
			Scan(&term.TermId, &term.Name, &term.StartDate, &term.EndDate)
	} else {
		tr.logger.Debug("Creating term without transaction")
		err = tr.db.Pool.QueryRow(ctx, sql, createTerm.Name, createTerm.StartDate, createTerm.EndDate).
			Scan(&term.TermId, &term.Name, &term.StartDate, &term.EndDate)
	}

	if err != nil {
		tr.logger.Error("Failed to create term", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &term, nil
}

func (tr *TermRepository) GetTermById(ctx context.Context, tx pgx.Tx, termId int32) (*Term, *exceptions.AppError) {
	if msg := domain_model.ValidateId(termId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid term ID", map[string]string{"term_id": msg})
	}

	sql := `
		SELECT term_id, name, start_date, end_date
		FROM terms
		WHERE term_id = $1;
	`

	var term Term
	var err error

	if tx != nil {
		tr.logger.Debug("Getting term by id in transaction", zap.Int32("term_id", termId))
		err = tx.QueryRow(ctx, sql, termId).
			Scan(&term.TermId, &term.Name, &term.StartDate, &term.EndDate)
	} else {
		tr.logger.Debug("Getting term by id without transaction", zap.Int32("term_id", termId))
		err = tr.db.Pool.QueryRow(ctx, sql, termId).
			Scan(&term.TermId, &term.Name, &term.StartDate, &term.EndDate)
	}

	if err != nil {
		tr.logger.Error("Failed to get term by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &term, nil
}

func (tr *TermRepository) GetTerms(ctx context.Context, tx pgx.Tx) ([]Term, *exceptions.AppError) {
	sql := `
		SELECT term_id, name, start_date, end_date
		FROM terms
		ORDER BY start_date, name;
	`

	var terms []Term
	var err error
	var rows pgx.Rows

	if tx != nil {
		tr.logger.Debug("Getting terms in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		tr.logger.Debug("Getting terms without transaction")
		rows, err = tr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		tr.logger.Error("Failed to get terms", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var term Term
		err = rows.Scan(&term.TermId, &term.Name, &term.StartDate, &term.EndDate)
		if err != nil {
			tr.logger.Error("Failed to scan term row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		terms = append(terms, term)
	}

	if err = rows.Err(); err != nil {
		tr.logger.Error("Error iterating terms rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return terms, nil
}

func (tr *TermRepository) DeleteTerm(ctx context.Context, tx pgx.Tx, termId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(termId); msg != "" {
		return exceptions.NewValidationError("Invalid term ID", map[string]string{"term_id": msg})
	}

	sql := "DELETE FROM terms WHERE term_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		tr.logger.Debug("Deleting term in transaction")
		cmdTag, err = tx.Exec(ctx, sql, termId)
	} else {
		tr.logger.Debug("Deleting term without transaction")
		cmdTag, err = tr.db.Pool.Exec(ctx, sql, termId)
	}
	if err != nil {
		tr.logger.Error("Failed to delete term", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		tr.logger.Error("Failed to delete term - term not found", zap.Int32("term_id", termId))
		return exceptions.NewNotFoundError("Term not found")
	}

	return nil
}


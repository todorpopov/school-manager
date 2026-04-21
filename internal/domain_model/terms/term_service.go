package terms

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ITermService interface {
	CreateTerm(ctx context.Context, tx pgx.Tx, createTerm *CreateTerm) (*Term, *exceptions.AppError)
	GetTermById(ctx context.Context, tx pgx.Tx, termId int32) (*Term, *exceptions.AppError)
	GetTerms(ctx context.Context, tx pgx.Tx) ([]Term, *exceptions.AppError)
	DeleteTerm(ctx context.Context, tx pgx.Tx, termId int32) *exceptions.AppError
}

type TermService struct {
	termRepo ITermRepository
}

func NewTermService(termRepo ITermRepository) *TermService {
	return &TermService{termRepo}
}

func (ts *TermService) CreateTerm(ctx context.Context, tx pgx.Tx, createTerm *CreateTerm) (*Term, *exceptions.AppError) {
	validationErr := ValidateCreateTerm(createTerm)
	if validationErr != nil {
		return nil, validationErr
	}

	return ts.termRepo.CreateTerm(ctx, tx, createTerm)
}

func (ts *TermService) GetTermById(ctx context.Context, tx pgx.Tx, termId int32) (*Term, *exceptions.AppError) {
	if msg := domain_model.ValidateId(termId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid term ID", map[string]string{"term_id": msg})
	}

	return ts.termRepo.GetTermById(ctx, tx, termId)
}

func (ts *TermService) GetTerms(ctx context.Context, tx pgx.Tx) ([]Term, *exceptions.AppError) {
	return ts.termRepo.GetTerms(ctx, tx)
}

func (ts *TermService) DeleteTerm(ctx context.Context, tx pgx.Tx, termId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(termId); msg != "" {
		return exceptions.NewValidationError("Invalid term ID", map[string]string{"term_id": msg})
	}

	return ts.termRepo.DeleteTerm(ctx, tx, termId)
}


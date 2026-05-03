package grades

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IGradeService interface {
	CreateGrade(ctx context.Context, tx pgx.Tx, createGrade *CreateGrade) (*Grade, *exceptions.AppError)
	BulkCreateGrades(ctx context.Context, tx pgx.Tx, payload *BulkCreateGrades) ([]Grade, *exceptions.AppError)
	GetGradeById(ctx context.Context, tx pgx.Tx, gradeId int32) (*Grade, *exceptions.AppError)
	GetGrades(ctx context.Context, tx pgx.Tx) ([]Grade, *exceptions.AppError)
	DeleteGrade(ctx context.Context, tx pgx.Tx, gradeId int32) *exceptions.AppError
}

type GradeService struct {
	gradeRepo IGradeRepository
}

func NewGradeService(gradeRepo IGradeRepository) *GradeService {
	return &GradeService{gradeRepo}
}

func (gs *GradeService) CreateGrade(ctx context.Context, tx pgx.Tx, createGrade *CreateGrade) (*Grade, *exceptions.AppError) {
	validationErr := ValidateCreateGrade(createGrade)
	if validationErr != nil {
		return nil, validationErr
	}

	return gs.gradeRepo.CreateGrade(ctx, tx, createGrade)
}

func (gs *GradeService) BulkCreateGrades(ctx context.Context, tx pgx.Tx, payload *BulkCreateGrades) ([]Grade, *exceptions.AppError) {
	if len(payload.Entries) == 0 {
		return nil, exceptions.NewValidationError("No entries provided", map[string]string{"entries": "must not be empty"})
	}
	for i, entry := range payload.Entries {
		if err := ValidateCreateGrade(&entry); err != nil {
			return nil, err
		}
		payload.Entries[i] = entry
	}
	return gs.gradeRepo.BulkCreateGrades(ctx, tx, payload.Entries)
}

func (gs *GradeService) GetGradeById(ctx context.Context, tx pgx.Tx, gradeId int32) (*Grade, *exceptions.AppError) {
	if msg := domain_model.ValidateId(gradeId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid grade ID", map[string]string{"grade_id": msg})
	}

	return gs.gradeRepo.GetGradeById(ctx, tx, gradeId)
}

func (gs *GradeService) GetGrades(ctx context.Context, tx pgx.Tx) ([]Grade, *exceptions.AppError) {
	return gs.gradeRepo.GetGrades(ctx, tx)
}

func (gs *GradeService) DeleteGrade(ctx context.Context, tx pgx.Tx, gradeId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(gradeId); msg != "" {
		return exceptions.NewValidationError("Invalid grade ID", map[string]string{"grade_id": msg})
	}

	return gs.gradeRepo.DeleteGrade(ctx, tx, gradeId)
}


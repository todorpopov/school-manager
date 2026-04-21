package subjects

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ISubjectService interface {
	CreateSubject(ctx context.Context, tx pgx.Tx, createSubject *CreateSubject) (*Subject, *exceptions.AppError)
	GetSubjectById(ctx context.Context, tx pgx.Tx, subjectId int32) (*Subject, *exceptions.AppError)
	GetSubjects(ctx context.Context, tx pgx.Tx) ([]Subject, *exceptions.AppError)
	DeleteSubject(ctx context.Context, tx pgx.Tx, subjectId int32) *exceptions.AppError
}

type SubjectService struct {
	subjectRepo ISubjectRepository
}

func NewSubjectService(subjectRepo ISubjectRepository) *SubjectService {
	return &SubjectService{subjectRepo}
}

func (ss *SubjectService) CreateSubject(ctx context.Context, tx pgx.Tx, createSubject *CreateSubject) (*Subject, *exceptions.AppError) {
	validationErr := ValidateCreateSubject(createSubject)
	if validationErr != nil {
		return nil, validationErr
	}

	return ss.subjectRepo.CreateSubject(ctx, tx, createSubject)
}

func (ss *SubjectService) GetSubjectById(ctx context.Context, tx pgx.Tx, subjectId int32) (*Subject, *exceptions.AppError) {
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	return ss.subjectRepo.GetSubjectById(ctx, tx, subjectId)
}

func (ss *SubjectService) GetSubjects(ctx context.Context, tx pgx.Tx) ([]Subject, *exceptions.AppError) {
	return ss.subjectRepo.GetSubjects(ctx, tx)
}

func (ss *SubjectService) DeleteSubject(ctx context.Context, tx pgx.Tx, subjectId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	return ss.subjectRepo.DeleteSubject(ctx, tx, subjectId)
}


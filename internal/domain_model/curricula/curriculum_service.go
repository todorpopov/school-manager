package curricula

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ICurriculumService interface {
	CreateCurriculum(ctx context.Context, tx pgx.Tx, createCurriculum *CreateCurriculum) (*Curriculum, *exceptions.AppError)
	GetCurriculumById(ctx context.Context, tx pgx.Tx, curriculumId int32) (*Curriculum, *exceptions.AppError)
	GetCurricula(ctx context.Context, tx pgx.Tx) ([]Curriculum, *exceptions.AppError)
	DeleteCurriculum(ctx context.Context, tx pgx.Tx, curriculumId int32) *exceptions.AppError
}

type CurriculumService struct {
	curriculumRepo ICurriculumRepository
}

func NewCurriculumService(curriculumRepo ICurriculumRepository) *CurriculumService {
	return &CurriculumService{curriculumRepo}
}

func (cs *CurriculumService) CreateCurriculum(ctx context.Context, tx pgx.Tx, createCurriculum *CreateCurriculum) (*Curriculum, *exceptions.AppError) {
	validationErr := ValidateCreateCurriculum(createCurriculum)
	if validationErr != nil {
		return nil, validationErr
	}

	return cs.curriculumRepo.CreateCurriculum(ctx, tx, createCurriculum)
}

func (cs *CurriculumService) GetCurriculumById(ctx context.Context, tx pgx.Tx, curriculumId int32) (*Curriculum, *exceptions.AppError) {
	if msg := domain_model.ValidateId(curriculumId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid curriculum ID", map[string]string{"curriculum_id": msg})
	}

	return cs.curriculumRepo.GetCurriculumById(ctx, tx, curriculumId)
}

func (cs *CurriculumService) GetCurricula(ctx context.Context, tx pgx.Tx) ([]Curriculum, *exceptions.AppError) {
	return cs.curriculumRepo.GetCurricula(ctx, tx)
}

func (cs *CurriculumService) DeleteCurriculum(ctx context.Context, tx pgx.Tx, curriculumId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(curriculumId); msg != "" {
		return exceptions.NewValidationError("Invalid curriculum ID", map[string]string{"curriculum_id": msg})
	}

	return cs.curriculumRepo.DeleteCurriculum(ctx, tx, curriculumId)
}


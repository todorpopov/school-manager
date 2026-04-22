package schools

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ISchoolService interface {
	CreateSchool(ctx context.Context, tx pgx.Tx, createSchool *CreateSchool) (*School, *exceptions.AppError)
	GetSchoolById(ctx context.Context, tx pgx.Tx, schoolId int32) (*School, *exceptions.AppError)
	GetSchools(ctx context.Context, tx pgx.Tx) ([]School, *exceptions.AppError)
	UpdateSchool(ctx context.Context, tx pgx.Tx, updateSchool *UpdateSchool) (*School, *exceptions.AppError)
	DeleteSchool(ctx context.Context, tx pgx.Tx, schoolId int32) *exceptions.AppError
}

type SchoolService struct {
	schoolRepo ISchoolRepository
}

func NewSchoolService(schoolRepo ISchoolRepository) *SchoolService {
	return &SchoolService{schoolRepo}
}

func (ss *SchoolService) CreateSchool(ctx context.Context, tx pgx.Tx, createSchool *CreateSchool) (*School, *exceptions.AppError) {
	validationErr := ValidateCreateSchool(createSchool)
	if validationErr != nil {
		return nil, validationErr
	}

	return ss.schoolRepo.CreateSchool(ctx, tx, createSchool)
}

func (ss *SchoolService) GetSchoolById(ctx context.Context, tx pgx.Tx, schoolId int32) (*School, *exceptions.AppError) {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	return ss.schoolRepo.GetSchoolById(ctx, tx, schoolId)
}

func (ss *SchoolService) GetSchools(ctx context.Context, tx pgx.Tx) ([]School, *exceptions.AppError) {
	return ss.schoolRepo.GetSchools(ctx, tx)
}

func (ss *SchoolService) UpdateSchool(ctx context.Context, tx pgx.Tx, updateSchool *UpdateSchool) (*School, *exceptions.AppError) {
	validationErr := ValidateUpdateSchool(updateSchool)
	if validationErr != nil {
		return nil, validationErr
	}

	return ss.schoolRepo.UpdateSchool(ctx, tx, updateSchool)
}

func (ss *SchoolService) DeleteSchool(ctx context.Context, tx pgx.Tx, schoolId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	return ss.schoolRepo.DeleteSchool(ctx, tx, schoolId)
}


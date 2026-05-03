package absences

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type IAbsenceService interface {
	CreateAbsence(ctx context.Context, tx pgx.Tx, createAbsence *CreateAbsence) (*Absence, *exceptions.AppError)
	BulkCreateAbsences(ctx context.Context, tx pgx.Tx, payload *BulkCreateAbsences) ([]Absence, *exceptions.AppError)
	GetAbsenceById(ctx context.Context, tx pgx.Tx, absenceId int32) (*Absence, *exceptions.AppError)
	GetAbsences(ctx context.Context, tx pgx.Tx) ([]Absence, *exceptions.AppError)
	DeleteAbsence(ctx context.Context, tx pgx.Tx, absenceId int32) *exceptions.AppError
}

type AbsenceService struct {
	absenceRepo IAbsenceRepository
}

func NewAbsenceService(absenceRepo IAbsenceRepository) *AbsenceService {
	return &AbsenceService{absenceRepo}
}

func (as *AbsenceService) CreateAbsence(ctx context.Context, tx pgx.Tx, createAbsence *CreateAbsence) (*Absence, *exceptions.AppError) {
	validationErr := ValidateCreateAbsence(createAbsence)
	if validationErr != nil {
		return nil, validationErr
	}

	return as.absenceRepo.CreateAbsence(ctx, tx, createAbsence)
}

func (as *AbsenceService) BulkCreateAbsences(ctx context.Context, tx pgx.Tx, payload *BulkCreateAbsences) ([]Absence, *exceptions.AppError) {
	if len(payload.Entries) == 0 {
		return nil, exceptions.NewValidationError("No entries provided", map[string]string{"entries": "must not be empty"})
	}
	for i, entry := range payload.Entries {
		if err := ValidateCreateAbsence(&entry); err != nil {
			return nil, err
		}
		payload.Entries[i] = entry
	}
	return as.absenceRepo.BulkCreateAbsences(ctx, tx, payload.Entries)
}

func (as *AbsenceService) GetAbsenceById(ctx context.Context, tx pgx.Tx, absenceId int32) (*Absence, *exceptions.AppError) {
	if msg := domain_model.ValidateId(absenceId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid absence ID", map[string]string{"absence_id": msg})
	}

	return as.absenceRepo.GetAbsenceById(ctx, tx, absenceId)
}

func (as *AbsenceService) GetAbsences(ctx context.Context, tx pgx.Tx) ([]Absence, *exceptions.AppError) {
	return as.absenceRepo.GetAbsences(ctx, tx)
}

func (as *AbsenceService) DeleteAbsence(ctx context.Context, tx pgx.Tx, absenceId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(absenceId); msg != "" {
		return exceptions.NewValidationError("Invalid absence ID", map[string]string{"absence_id": msg})
	}

	return as.absenceRepo.DeleteAbsence(ctx, tx, absenceId)
}


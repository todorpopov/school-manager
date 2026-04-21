package teachers

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type ITeacherSubjectService interface {
	LinkSubjectToTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError
	UnlinkSubjectFromTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError
	GetSubjectsForTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) ([]subjects.Subject, *exceptions.AppError)
	GetTeachersForSubject(ctx context.Context, tx pgx.Tx, subjectId int32) ([]Teacher, *exceptions.AppError)
	CanTeacherTeachSubject(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) (bool, *exceptions.AppError)
}

type TeacherSubjectService struct {
	teacherSubjectRepo ITeacherSubjectRepository
}

func NewTeacherSubjectService(teacherSubjectRepo ITeacherSubjectRepository) *TeacherSubjectService {
	return &TeacherSubjectService{teacherSubjectRepo}
}

func (tss *TeacherSubjectService) LinkSubjectToTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	return tss.teacherSubjectRepo.LinkSubjectToTeacher(ctx, tx, teacherId, subjectId)
}

func (tss *TeacherSubjectService) UnlinkSubjectFromTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	return tss.teacherSubjectRepo.UnlinkSubjectFromTeacher(ctx, tx, teacherId, subjectId)
}

func (tss *TeacherSubjectService) GetSubjectsForTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) ([]subjects.Subject, *exceptions.AppError) {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}

	return tss.teacherSubjectRepo.GetSubjectsForTeacher(ctx, tx, teacherId)
}

func (tss *TeacherSubjectService) GetTeachersForSubject(ctx context.Context, tx pgx.Tx, subjectId int32) ([]Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	return tss.teacherSubjectRepo.GetTeachersForSubject(ctx, tx, subjectId)
}

func (tss *TeacherSubjectService) CanTeacherTeachSubject(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) (bool, *exceptions.AppError) {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return false, exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return false, exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	return tss.teacherSubjectRepo.CanTeacherTeachSubject(ctx, tx, teacherId, subjectId)
}


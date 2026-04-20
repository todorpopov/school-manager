package teachers

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ITeacherSubjectRepository interface {
	LinkSubjectToTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError
	UnlinkSubjectFromTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError
	GetSubjectsForTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) ([]subjects.Subject, *exceptions.AppError)
	GetTeachersForSubject(ctx context.Context, tx pgx.Tx, subjectId int32) ([]Teacher, *exceptions.AppError)
	CanTeacherTeachSubject(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) (bool, *exceptions.AppError)
}

type TeacherSubjectRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewTeacherSubjectRepository(db *persistence.Database, logger *zap.Logger) *TeacherSubjectRepository {
	return &TeacherSubjectRepository{db, logger}
}

func (tsr *TeacherSubjectRepository) LinkSubjectToTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	sqlStmt := `
		INSERT INTO teacher_subjects (teacher_id, subject_id)
		VALUES ($1, $2);
	`

	var err error
	if tx != nil {
		tsr.logger.Debug("Linking subject to teacher in transaction", zap.Int32("teacher_id", teacherId), zap.Int32("subject_id", subjectId))
		_, err = tx.Exec(ctx, sqlStmt, teacherId, subjectId)
	} else {
		tsr.logger.Debug("Linking subject to teacher without transaction", zap.Int32("teacher_id", teacherId), zap.Int32("subject_id", subjectId))
		_, err = tsr.db.Pool.Exec(ctx, sqlStmt, teacherId, subjectId)
	}

	if err != nil {
		tsr.logger.Error("Failed to link subject to teacher", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	return nil
}

func (tsr *TeacherSubjectRepository) UnlinkSubjectFromTeacher(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	sqlStmt := `
		DELETE FROM teacher_subjects
		WHERE teacher_id = $1 AND subject_id = $2;
	`

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		tsr.logger.Debug("Unlinking subject from teacher in transaction", zap.Int32("teacher_id", teacherId), zap.Int32("subject_id", subjectId))
		cmdTag, err = tx.Exec(ctx, sqlStmt, teacherId, subjectId)
	} else {
		tsr.logger.Debug("Unlinking subject from teacher without transaction", zap.Int32("teacher_id", teacherId), zap.Int32("subject_id", subjectId))
		cmdTag, err = tsr.db.Pool.Exec(ctx, sqlStmt, teacherId, subjectId)
	}

	if err != nil {
		tsr.logger.Error("Failed to unlink subject from teacher", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		tsr.logger.Error("Failed to unlink subject from teacher - relationship not found")
		return exceptions.NewNotFoundError("Teacher-subject relationship not found")
	}

	return nil
}

func (tsr *TeacherSubjectRepository) GetSubjectsForTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) ([]subjects.Subject, *exceptions.AppError) {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}

	sqlStmt := `
		SELECT s.subject_id, s.subject_name
		FROM teacher_subjects ts
		INNER JOIN subjects s ON ts.subject_id = s.subject_id
		WHERE ts.teacher_id = $1
		ORDER BY s.subject_name;
	`

	var teacherSubjects []subjects.Subject
	var err error
	var rows pgx.Rows

	if tx != nil {
		tsr.logger.Debug("Getting subjects for teacher in transaction", zap.Int32("teacher_id", teacherId))
		rows, err = tx.Query(ctx, sqlStmt, teacherId)
	} else {
		tsr.logger.Debug("Getting subjects for teacher without transaction", zap.Int32("teacher_id", teacherId))
		rows, err = tsr.db.Pool.Query(ctx, sqlStmt, teacherId)
	}

	if err != nil {
		tsr.logger.Error("Failed to get subjects for teacher", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var subject subjects.Subject
		err = rows.Scan(&subject.SubjectId, &subject.SubjectName)
		if err != nil {
			tsr.logger.Error("Failed to scan subject row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		teacherSubjects = append(teacherSubjects, subject)
	}

	if err = rows.Err(); err != nil {
		tsr.logger.Error("Error iterating subject rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return teacherSubjects, nil
}

func (tsr *TeacherSubjectRepository) GetTeachersForSubject(ctx context.Context, tx pgx.Tx, subjectId int32) ([]Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	sqlStmt := `
		SELECT t.teacher_id, t.user_id, u.first_name, u.last_name, u.email
		FROM teacher_subjects ts
		INNER JOIN teachers t ON ts.teacher_id = t.teacher_id
		INNER JOIN users u ON t.user_id = u.user_id
		WHERE ts.subject_id = $1
		ORDER BY u.last_name, u.first_name;
	`

	var subjectTeachers []Teacher
	var err error
	var rows pgx.Rows

	if tx != nil {
		tsr.logger.Debug("Getting teachers for subject in transaction", zap.Int32("subject_id", subjectId))
		rows, err = tx.Query(ctx, sqlStmt, subjectId)
	} else {
		tsr.logger.Debug("Getting teachers for subject without transaction", zap.Int32("subject_id", subjectId))
		rows, err = tsr.db.Pool.Query(ctx, sqlStmt, subjectId)
	}

	if err != nil {
		tsr.logger.Error("Failed to get teachers for subject", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var teacher Teacher
		err = rows.Scan(&teacher.TeacherId, &teacher.UserId, &teacher.FirstName, &teacher.LastName, &teacher.Email)
		if err != nil {
			tsr.logger.Error("Failed to scan teacher row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		subjectTeachers = append(subjectTeachers, teacher)
	}

	if err = rows.Err(); err != nil {
		tsr.logger.Error("Error iterating teacher rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return subjectTeachers, nil
}

func (tsr *TeacherSubjectRepository) CanTeacherTeachSubject(ctx context.Context, tx pgx.Tx, teacherId int32, subjectId int32) (bool, *exceptions.AppError) {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return false, exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return false, exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	sqlStmt := `
		SELECT EXISTS(
			SELECT 1
			FROM teacher_subjects
			WHERE teacher_id = $1 AND subject_id = $2
		);
	`

	var exists bool
	var err error

	if tx != nil {
		tsr.logger.Debug("Checking if teacher can teach subject in transaction", zap.Int32("teacher_id", teacherId), zap.Int32("subject_id", subjectId))
		err = tx.QueryRow(ctx, sqlStmt, teacherId, subjectId).Scan(&exists)
	} else {
		tsr.logger.Debug("Checking if teacher can teach subject without transaction", zap.Int32("teacher_id", teacherId), zap.Int32("subject_id", subjectId))
		err = tsr.db.Pool.QueryRow(ctx, sqlStmt, teacherId, subjectId).Scan(&exists)
	}

	if err != nil {
		tsr.logger.Error("Failed to check if teacher can teach subject", zap.Error(err))
		return false, exceptions.PgErrorToAppError(err)
	}

	return exists, nil
}


package teachers

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ITeacherRepository interface {
	CreateTeacher(ctx context.Context, tx pgx.Tx, userId int32) (*Teacher, *exceptions.AppError)
	GetTeacherById(ctx context.Context, tx pgx.Tx, teacherId int32) (*Teacher, *exceptions.AppError)
	GetTeacherByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Teacher, *exceptions.AppError)
	GetTeachers(ctx context.Context, tx pgx.Tx) ([]Teacher, *exceptions.AppError)
	DeleteTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) *exceptions.AppError
}

type TeacherRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewTeacherRepository(db *persistence.Database, logger *zap.Logger) *TeacherRepository {
	return &TeacherRepository{db, logger}
}

func (tr *TeacherRepository) CreateTeacher(ctx context.Context, tx pgx.Tx, userId int32) (*Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		INSERT INTO teachers (user_id)
		VALUES ($1)
		RETURNING teacher_id, user_id;
	`

	var teacher Teacher
	var err error

	if tx != nil {
		tr.logger.Debug("Creating teacher in transaction")
		err = tx.QueryRow(ctx, sql, userId).
			Scan(&teacher.TeacherId, &teacher.UserId)
	} else {
		tr.logger.Debug("Creating teacher without transaction")
		err = tr.db.Pool.QueryRow(ctx, sql, userId).
			Scan(&teacher.TeacherId, &teacher.UserId)
	}

	if err != nil {
		tr.logger.Error("Failed to create teacher", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &teacher, nil
}

func (tr *TeacherRepository) GetTeacherById(ctx context.Context, tx pgx.Tx, teacherId int32) (*Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}

	sql := `
		SELECT t.teacher_id, t.user_id, u.first_name, u.last_name, u.email
		FROM teachers t
		INNER JOIN users u ON t.user_id = u.user_id
		WHERE t.teacher_id = $1;
	`

	var teacher Teacher
	var err error

	if tx != nil {
		tr.logger.Debug("Getting teacher by id in transaction", zap.Int32("teacher_id", teacherId))
		err = tx.QueryRow(ctx, sql, teacherId).
			Scan(&teacher.TeacherId, &teacher.UserId, &teacher.FirstName, &teacher.LastName, &teacher.Email)
	} else {
		tr.logger.Debug("Getting teacher by id without transaction", zap.Int32("teacher_id", teacherId))
		err = tr.db.Pool.QueryRow(ctx, sql, teacherId).
			Scan(&teacher.TeacherId, &teacher.UserId, &teacher.FirstName, &teacher.LastName, &teacher.Email)
	}

	if err != nil {
		tr.logger.Error("Failed to get teacher by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &teacher, nil
}

func (tr *TeacherRepository) GetTeacherByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Teacher, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		SELECT t.teacher_id, t.user_id, u.first_name, u.last_name, u.email
		FROM teachers t
		INNER JOIN users u ON t.user_id = u.user_id
		WHERE t.user_id = $1;
	`

	var teacher Teacher
	var err error

	if tx != nil {
		tr.logger.Debug("Getting teacher by user_id in transaction", zap.Int32("user_id", userId))
		err = tx.QueryRow(ctx, sql, userId).
			Scan(&teacher.TeacherId, &teacher.UserId, &teacher.FirstName, &teacher.LastName, &teacher.Email)
	} else {
		tr.logger.Debug("Getting teacher by user_id without transaction", zap.Int32("user_id", userId))
		err = tr.db.Pool.QueryRow(ctx, sql, userId).
			Scan(&teacher.TeacherId, &teacher.UserId, &teacher.FirstName, &teacher.LastName, &teacher.Email)
	}

	if err != nil {
		tr.logger.Error("Failed to get teacher by user_id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &teacher, nil
}

func (tr *TeacherRepository) GetTeachers(ctx context.Context, tx pgx.Tx) ([]Teacher, *exceptions.AppError) {
	sql := `
		SELECT t.teacher_id, t.user_id, u.first_name, u.last_name, u.email
		FROM teachers t
		INNER JOIN users u ON t.user_id = u.user_id;
	`

	var teachers []Teacher
	var err error
	var rows pgx.Rows

	if tx != nil {
		tr.logger.Debug("Getting teachers in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		tr.logger.Debug("Getting teachers without transaction")
		rows, err = tr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		tr.logger.Error("Failed to get teachers", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var teacher Teacher
		err = rows.Scan(&teacher.TeacherId, &teacher.UserId, &teacher.FirstName, &teacher.LastName, &teacher.Email)
		if err != nil {
			tr.logger.Error("Failed to scan teacher row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		tr.logger.Error("Error iterating teachers rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return teachers, nil
}


func (tr *TeacherRepository) DeleteTeacher(ctx context.Context, tx pgx.Tx, teacherId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(teacherId); msg != "" {
		return exceptions.NewValidationError("Invalid teacher ID", map[string]string{"teacher_id": msg})
	}

	sql := "DELETE FROM teachers WHERE teacher_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		tr.logger.Debug("Deleting teacher in transaction")
		cmdTag, err = tx.Exec(ctx, sql, teacherId)
	} else {
		tr.logger.Debug("Deleting teacher without transaction")
		cmdTag, err = tr.db.Pool.Exec(ctx, sql, teacherId)
	}
	if err != nil {
		tr.logger.Error("Failed to delete teacher", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		tr.logger.Error("Failed to delete teacher - teacher not found", zap.Int32("teacher_id", teacherId))
		return exceptions.NewNotFoundError("Teacher not found")
	}

	return nil
}


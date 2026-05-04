package classes

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IClassRepository interface {
	CreateClass(ctx context.Context, tx pgx.Tx, createClass *CreateClass) (*Class, *exceptions.AppError)
	GetClassById(ctx context.Context, tx pgx.Tx, classId int32) (*Class, *exceptions.AppError)
	GetClasses(ctx context.Context, tx pgx.Tx) ([]Class, *exceptions.AppError)
	GetClassesBySchoolId(ctx context.Context, tx pgx.Tx, schoolId int32) ([]Class, *exceptions.AppError)
	DeleteClass(ctx context.Context, tx pgx.Tx, classId int32) *exceptions.AppError
}

type ClassRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewClassRepository(db *persistence.Database, logger *zap.Logger) *ClassRepository {
	return &ClassRepository{db, logger}
}

func (cr *ClassRepository) CreateClass(ctx context.Context, tx pgx.Tx, createClass *CreateClass) (*Class, *exceptions.AppError) {
	sql := `
		INSERT INTO classes (school_id, grade_level, class_name)
		VALUES ($1, $2, $3)
		RETURNING class_id, school_id, grade_level, class_name;
	`

	var class Class
	var err error

	if tx != nil {
		cr.logger.Debug("Creating class in transaction")
		err = tx.QueryRow(ctx, sql, createClass.SchoolId, createClass.GradeLevel, createClass.ClassName).
			Scan(&class.ClassId, &class.SchoolId, &class.GradeLevel, &class.ClassName)
	} else {
		cr.logger.Debug("Creating class without transaction")
		err = cr.db.Pool.QueryRow(ctx, sql, createClass.SchoolId, createClass.GradeLevel, createClass.ClassName).
			Scan(&class.ClassId, &class.SchoolId, &class.GradeLevel, &class.ClassName)
	}

	if err != nil {
		cr.logger.Error("Failed to create class", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &class, nil
}

func (cr *ClassRepository) GetClassById(ctx context.Context, tx pgx.Tx, classId int32) (*Class, *exceptions.AppError) {
	if msg := domain_model.ValidateId(classId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid class ID", map[string]string{"class_id": msg})
	}

	sql := `
		SELECT class_id, school_id, grade_level, class_name
		FROM classes
		WHERE class_id = $1;
	`

	var class Class
	var err error

	if tx != nil {
		cr.logger.Debug("Getting class by id in transaction", zap.Int32("class_id", classId))
		err = tx.QueryRow(ctx, sql, classId).
			Scan(&class.ClassId, &class.SchoolId, &class.GradeLevel, &class.ClassName)
	} else {
		cr.logger.Debug("Getting class by id without transaction", zap.Int32("class_id", classId))
		err = cr.db.Pool.QueryRow(ctx, sql, classId).
			Scan(&class.ClassId, &class.SchoolId, &class.GradeLevel, &class.ClassName)
	}

	if err != nil {
		cr.logger.Error("Failed to get class by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &class, nil
}

func (cr *ClassRepository) GetClasses(ctx context.Context, tx pgx.Tx) ([]Class, *exceptions.AppError) {
	sql := `
		SELECT class_id, school_id, grade_level, class_name
		FROM classes
		ORDER BY school_id, grade_level, class_name;
	`

	var classes []Class
	var err error
	var rows pgx.Rows

	if tx != nil {
		cr.logger.Debug("Getting classes in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		cr.logger.Debug("Getting classes without transaction")
		rows, err = cr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		cr.logger.Error("Failed to get classes", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var class Class
		err = rows.Scan(&class.ClassId, &class.SchoolId, &class.GradeLevel, &class.ClassName)
		if err != nil {
			cr.logger.Error("Failed to scan class row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		cr.logger.Error("Error iterating classes rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return classes, nil
}

func (cr *ClassRepository) GetClassesBySchoolId(ctx context.Context, tx pgx.Tx, schoolId int32) ([]Class, *exceptions.AppError) {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	sql := `
		SELECT class_id, school_id, grade_level, class_name
		FROM classes
		WHERE school_id = $1
		ORDER BY grade_level, class_name;
	`

	var classes []Class
	var err error
	var rows pgx.Rows

	if tx != nil {
		cr.logger.Debug("Getting classes by school id in transaction", zap.Int32("school_id", schoolId))
		rows, err = tx.Query(ctx, sql, schoolId)
	} else {
		cr.logger.Debug("Getting classes by school id without transaction", zap.Int32("school_id", schoolId))
		rows, err = cr.db.Pool.Query(ctx, sql, schoolId)
	}

	if err != nil {
		cr.logger.Error("Failed to get classes by school id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var class Class
		err = rows.Scan(&class.ClassId, &class.SchoolId, &class.GradeLevel, &class.ClassName)
		if err != nil {
			cr.logger.Error("Failed to scan class row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		classes = append(classes, class)
	}

	if err = rows.Err(); err != nil {
		cr.logger.Error("Error iterating classes rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return classes, nil
}

func (cr *ClassRepository) DeleteClass(ctx context.Context, tx pgx.Tx, classId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(classId); msg != "" {
		return exceptions.NewValidationError("Invalid class ID", map[string]string{"class_id": msg})
	}

	sql := "DELETE FROM classes WHERE class_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		cr.logger.Debug("Deleting class in transaction")
		cmdTag, err = tx.Exec(ctx, sql, classId)
	} else {
		cr.logger.Debug("Deleting class without transaction")
		cmdTag, err = cr.db.Pool.Exec(ctx, sql, classId)
	}
	if err != nil {
		cr.logger.Error("Failed to delete class", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		cr.logger.Error("Failed to delete class - class not found", zap.Int32("class_id", classId))
		return exceptions.NewNotFoundError("Class not found")
	}

	return nil
}

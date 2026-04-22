package schools

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ISchoolRepository interface {
	CreateSchool(ctx context.Context, tx pgx.Tx, createSchool *CreateSchool) (*School, *exceptions.AppError)
	GetSchoolById(ctx context.Context, tx pgx.Tx, schoolId int32) (*School, *exceptions.AppError)
	GetSchools(ctx context.Context, tx pgx.Tx) ([]School, *exceptions.AppError)
	UpdateSchool(ctx context.Context, tx pgx.Tx, updateSchool *UpdateSchool) (*School, *exceptions.AppError)
	DeleteSchool(ctx context.Context, tx pgx.Tx, schoolId int32) *exceptions.AppError
}

type SchoolRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewSchoolRepository(db *persistence.Database, logger *zap.Logger) *SchoolRepository {
	return &SchoolRepository{db, logger}
}

func (sr *SchoolRepository) CreateSchool(ctx context.Context, tx pgx.Tx, createSchool *CreateSchool) (*School, *exceptions.AppError) {
	sql := `
		INSERT INTO schools (school_name, school_address)
		VALUES ($1, $2)
		RETURNING school_id, school_name, school_address;
	`

	var school School
	var err error

	if tx != nil {
		sr.logger.Debug("Creating school in transaction")
		err = tx.QueryRow(ctx, sql, createSchool.SchoolName, createSchool.SchoolAddress).
			Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
	} else {
		sr.logger.Debug("Creating school without transaction")
		err = sr.db.Pool.QueryRow(ctx, sql, createSchool.SchoolName, createSchool.SchoolAddress).
			Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
	}

	if err != nil {
		sr.logger.Error("Failed to create school", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &school, nil
}

func (sr *SchoolRepository) GetSchoolById(ctx context.Context, tx pgx.Tx, schoolId int32) (*School, *exceptions.AppError) {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	sql := `
		SELECT school_id, school_name, school_address
		FROM schools
		WHERE school_id = $1;
	`

	var school School
	var err error

	if tx != nil {
		sr.logger.Debug("Getting school by id in transaction", zap.Int32("school_id", schoolId))
		err = tx.QueryRow(ctx, sql, schoolId).
			Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
	} else {
		sr.logger.Debug("Getting school by id without transaction", zap.Int32("school_id", schoolId))
		err = sr.db.Pool.QueryRow(ctx, sql, schoolId).
			Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
	}

	if err != nil {
		sr.logger.Error("Failed to get school by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &school, nil
}

func (sr *SchoolRepository) GetSchools(ctx context.Context, tx pgx.Tx) ([]School, *exceptions.AppError) {
	sql := `
		SELECT school_id, school_name, school_address
		FROM schools
		ORDER BY school_name;
	`

	var schools []School
	var err error
	var rows pgx.Rows

	if tx != nil {
		sr.logger.Debug("Getting schools in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		sr.logger.Debug("Getting schools without transaction")
		rows, err = sr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		sr.logger.Error("Failed to get schools", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var school School
		err = rows.Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
		if err != nil {
			sr.logger.Error("Failed to scan school row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		schools = append(schools, school)
	}

	if err = rows.Err(); err != nil {
		sr.logger.Error("Error iterating schools rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return schools, nil
}

func (sr *SchoolRepository) UpdateSchool(ctx context.Context, tx pgx.Tx, updateSchool *UpdateSchool) (*School, *exceptions.AppError) {
	if msg := domain_model.ValidateId(updateSchool.SchoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	sql := `
		UPDATE schools
		SET school_name = $2, school_address = $3
		WHERE school_id = $1
		RETURNING school_id, school_name, school_address;
	`

	var school School
	var err error

	if tx != nil {
		sr.logger.Debug("Updating school in transaction", zap.Int32("school_id", updateSchool.SchoolId))
		err = tx.QueryRow(ctx, sql, updateSchool.SchoolId, updateSchool.SchoolName, updateSchool.SchoolAddress).
			Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
	} else {
		sr.logger.Debug("Updating school without transaction", zap.Int32("school_id", updateSchool.SchoolId))
		err = sr.db.Pool.QueryRow(ctx, sql, updateSchool.SchoolId, updateSchool.SchoolName, updateSchool.SchoolAddress).
			Scan(&school.SchoolId, &school.SchoolName, &school.SchoolAddress)
	}

	if err != nil {
		sr.logger.Error("Failed to update school", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &school, nil
}

func (sr *SchoolRepository) DeleteSchool(ctx context.Context, tx pgx.Tx, schoolId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	sql := "DELETE FROM schools WHERE school_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		sr.logger.Debug("Deleting school in transaction")
		cmdTag, err = tx.Exec(ctx, sql, schoolId)
	} else {
		sr.logger.Debug("Deleting school without transaction")
		cmdTag, err = sr.db.Pool.Exec(ctx, sql, schoolId)
	}
	if err != nil {
		sr.logger.Error("Failed to delete school", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		sr.logger.Error("Failed to delete school - school not found", zap.Int32("school_id", schoolId))
		return exceptions.NewNotFoundError("School not found")
	}

	return nil
}


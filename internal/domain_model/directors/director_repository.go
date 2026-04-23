package directors

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IDirectorRepository interface {
	CreateDirector(ctx context.Context, tx pgx.Tx, schoolId int32, userId int32) (*Director, *exceptions.AppError)
	GetDirectorById(ctx context.Context, tx pgx.Tx, directorId int32) (*Director, *exceptions.AppError)
	GetDirectorByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Director, *exceptions.AppError)
	GetDirectors(ctx context.Context, tx pgx.Tx) ([]Director, *exceptions.AppError)
	DeleteDirector(ctx context.Context, tx pgx.Tx, directorId int32) *exceptions.AppError
}

type DirectorRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewDirectorRepository(db *persistence.Database, logger *zap.Logger) *DirectorRepository {
	return &DirectorRepository{db, logger}
}

func (dr *DirectorRepository) CreateDirector(ctx context.Context, tx pgx.Tx, schoolId int32, userId int32) (*Director, *exceptions.AppError) {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}

	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		INSERT INTO directors (school_id, user_id)
		VALUES ($1, $2)
		RETURNING director_id;
	`

	var directorId int32
	var err error

	if tx != nil {
		dr.logger.Debug("Creating director in transaction")
		err = tx.QueryRow(ctx, sql, schoolId, userId).Scan(&directorId)
	} else {
		dr.logger.Debug("Creating director without transaction")
		err = dr.db.Pool.QueryRow(ctx, sql, schoolId, userId).Scan(&directorId)
	}

	if err != nil {
		dr.logger.Error("Failed to create director", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return dr.GetDirectorById(ctx, tx, directorId)
}

func (dr *DirectorRepository) GetDirectorById(ctx context.Context, tx pgx.Tx, directorId int32) (*Director, *exceptions.AppError) {
	if msg := domain_model.ValidateId(directorId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid director ID", map[string]string{"director_id": msg})
	}

	sql := `
		SELECT d.director_id, d.user_id, u.first_name, u.last_name, u.email,
		       s.school_id, s.school_name, s.school_address
		FROM directors d
		INNER JOIN users u ON d.user_id = u.user_id
		INNER JOIN schools s ON d.school_id = s.school_id
		WHERE d.director_id = $1;
	`

	var director Director
	var schoolId int32
	var schoolName, schoolAddress string
	var err error

	if tx != nil {
		dr.logger.Debug("Getting director by id in transaction", zap.Int32("director_id", directorId))
		err = tx.QueryRow(ctx, sql, directorId).
			Scan(&director.DirectorId, &director.UserId, &director.FirstName, &director.LastName, &director.Email,
				&schoolId, &schoolName, &schoolAddress)
	} else {
		dr.logger.Debug("Getting director by id without transaction", zap.Int32("director_id", directorId))
		err = dr.db.Pool.QueryRow(ctx, sql, directorId).
			Scan(&director.DirectorId, &director.UserId, &director.FirstName, &director.LastName, &director.Email,
				&schoolId, &schoolName, &schoolAddress)
	}

	if err != nil {
		dr.logger.Error("Failed to get director by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	director.School = &schools.School{
		SchoolId:      schoolId,
		SchoolName:    schoolName,
		SchoolAddress: schoolAddress,
	}

	return &director, nil
}

func (dr *DirectorRepository) GetDirectorByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Director, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		SELECT d.director_id, d.user_id, u.first_name, u.last_name, u.email,
		       s.school_id, s.school_name, s.school_address
		FROM directors d
		INNER JOIN users u ON d.user_id = u.user_id
		INNER JOIN schools s ON d.school_id = s.school_id
		WHERE d.user_id = $1;
	`

	var director Director
	var schoolId int32
	var schoolName, schoolAddress string
	var err error

	if tx != nil {
		dr.logger.Debug("Getting director by user_id in transaction", zap.Int32("user_id", userId))
		err = tx.QueryRow(ctx, sql, userId).
			Scan(&director.DirectorId, &director.UserId, &director.FirstName, &director.LastName, &director.Email,
				&schoolId, &schoolName, &schoolAddress)
	} else {
		dr.logger.Debug("Getting director by user_id without transaction", zap.Int32("user_id", userId))
		err = dr.db.Pool.QueryRow(ctx, sql, userId).
			Scan(&director.DirectorId, &director.UserId, &director.FirstName, &director.LastName, &director.Email,
				&schoolId, &schoolName, &schoolAddress)
	}

	if err != nil {
		dr.logger.Error("Failed to get director by user_id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	director.School = &schools.School{
		SchoolId:      schoolId,
		SchoolName:    schoolName,
		SchoolAddress: schoolAddress,
	}

	return &director, nil
}

func (dr *DirectorRepository) GetDirectors(ctx context.Context, tx pgx.Tx) ([]Director, *exceptions.AppError) {
	sql := `
		SELECT d.director_id, d.user_id, u.first_name, u.last_name, u.email,
		       s.school_id, s.school_name, s.school_address
		FROM directors d
		INNER JOIN users u ON d.user_id = u.user_id
		INNER JOIN schools s ON d.school_id = s.school_id;
	`

	var directors []Director
	var err error
	var rows pgx.Rows

	if tx != nil {
		dr.logger.Debug("Getting directors in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		dr.logger.Debug("Getting directors without transaction")
		rows, err = dr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		dr.logger.Error("Failed to get directors", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var director Director
		var schoolId int32
		var schoolName, schoolAddress string

		err = rows.Scan(&director.DirectorId, &director.UserId, &director.FirstName, &director.LastName, &director.Email,
			&schoolId, &schoolName, &schoolAddress)
		if err != nil {
			dr.logger.Error("Failed to scan director row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		director.School = &schools.School{
			SchoolId:      schoolId,
			SchoolName:    schoolName,
			SchoolAddress: schoolAddress,
		}

		directors = append(directors, director)
	}

	if err = rows.Err(); err != nil {
		dr.logger.Error("Error iterating directors rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return directors, nil
}

func (dr *DirectorRepository) DeleteDirector(ctx context.Context, tx pgx.Tx, directorId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(directorId); msg != "" {
		return exceptions.NewValidationError("Invalid director ID", map[string]string{"director_id": msg})
	}

	sql := "DELETE FROM directors WHERE director_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		dr.logger.Debug("Deleting director in transaction")
		cmdTag, err = tx.Exec(ctx, sql, directorId)
	} else {
		dr.logger.Debug("Deleting director without transaction")
		cmdTag, err = dr.db.Pool.Exec(ctx, sql, directorId)
	}
	if err != nil {
		dr.logger.Error("Failed to delete director", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		dr.logger.Error("Failed to delete director - director not found", zap.Int32("director_id", directorId))
		return exceptions.NewNotFoundError("Director not found")
	}

	return nil
}

package subjects

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ISubjectRepository interface {
	CreateSubject(ctx context.Context, tx pgx.Tx, createSubject *CreateSubject) (*Subject, *exceptions.AppError)
	GetSubjectById(ctx context.Context, tx pgx.Tx, subjectId int32) (*Subject, *exceptions.AppError)
	GetSubjects(ctx context.Context, tx pgx.Tx) ([]Subject, *exceptions.AppError)
	DeleteSubject(ctx context.Context, tx pgx.Tx, subjectId int32) *exceptions.AppError
}

type SubjectRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewSubjectRepository(db *persistence.Database, logger *zap.Logger) *SubjectRepository {
	return &SubjectRepository{db, logger}
}

func (sr *SubjectRepository) CreateSubject(ctx context.Context, tx pgx.Tx, createSubject *CreateSubject) (*Subject, *exceptions.AppError) {
	sql := `
		INSERT INTO subjects (subject_name)
		VALUES ($1)
		RETURNING subject_id, subject_name;
	`

	var subject Subject
	var err error

	if tx != nil {
		sr.logger.Debug("Creating subject in transaction")
		err = tx.QueryRow(ctx, sql, createSubject.SubjectName).
			Scan(&subject.SubjectId, &subject.SubjectName)
	} else {
		sr.logger.Debug("Creating subject without transaction")
		err = sr.db.Pool.QueryRow(ctx, sql, createSubject.SubjectName).
			Scan(&subject.SubjectId, &subject.SubjectName)
	}

	if err != nil {
		sr.logger.Error("Failed to create subject", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &subject, nil
}

func (sr *SubjectRepository) GetSubjectById(ctx context.Context, tx pgx.Tx, subjectId int32) (*Subject, *exceptions.AppError) {
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	sql := `
		SELECT subject_id, subject_name
		FROM subjects
		WHERE subject_id = $1;
	`

	var subject Subject
	var err error

	if tx != nil {
		sr.logger.Debug("Getting subject by id in transaction", zap.Int32("subject_id", subjectId))
		err = tx.QueryRow(ctx, sql, subjectId).
			Scan(&subject.SubjectId, &subject.SubjectName)
	} else {
		sr.logger.Debug("Getting subject by id without transaction", zap.Int32("subject_id", subjectId))
		err = sr.db.Pool.QueryRow(ctx, sql, subjectId).
			Scan(&subject.SubjectId, &subject.SubjectName)
	}

	if err != nil {
		sr.logger.Error("Failed to get subject by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return &subject, nil
}

func (sr *SubjectRepository) GetSubjects(ctx context.Context, tx pgx.Tx) ([]Subject, *exceptions.AppError) {
	sql := `
		SELECT subject_id, subject_name
		FROM subjects
		ORDER BY subject_name;
	`

	var subjects []Subject
	var err error
	var rows pgx.Rows

	if tx != nil {
		sr.logger.Debug("Getting subjects in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		sr.logger.Debug("Getting subjects without transaction")
		rows, err = sr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		sr.logger.Error("Failed to get subjects", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var subject Subject
		err = rows.Scan(&subject.SubjectId, &subject.SubjectName)
		if err != nil {
			sr.logger.Error("Failed to scan subject row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		subjects = append(subjects, subject)
	}

	if err = rows.Err(); err != nil {
		sr.logger.Error("Error iterating subjects rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return subjects, nil
}

func (sr *SubjectRepository) DeleteSubject(ctx context.Context, tx pgx.Tx, subjectId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(subjectId); msg != "" {
		return exceptions.NewValidationError("Invalid subject ID", map[string]string{"subject_id": msg})
	}

	sql := "DELETE FROM subjects WHERE subject_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		sr.logger.Debug("Deleting subject in transaction")
		cmdTag, err = tx.Exec(ctx, sql, subjectId)
	} else {
		sr.logger.Debug("Deleting subject without transaction")
		cmdTag, err = sr.db.Pool.Exec(ctx, sql, subjectId)
	}
	if err != nil {
		sr.logger.Error("Failed to delete subject", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		sr.logger.Error("Failed to delete subject - subject not found", zap.Int32("subject_id", subjectId))
		return exceptions.NewNotFoundError("Subject not found")
	}

	return nil
}


package curricula

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/domain_model/terms"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type ICurriculumRepository interface {
	CreateCurriculum(ctx context.Context, tx pgx.Tx, createCurriculum *CreateCurriculum) (*Curriculum, *exceptions.AppError)
	GetCurriculumById(ctx context.Context, tx pgx.Tx, curriculumId int32) (*Curriculum, *exceptions.AppError)
	GetCurricula(ctx context.Context, tx pgx.Tx) ([]Curriculum, *exceptions.AppError)
	DeleteCurriculum(ctx context.Context, tx pgx.Tx, curriculumId int32) *exceptions.AppError
}

type CurriculumRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewCurriculumRepository(db *persistence.Database, logger *zap.Logger) *CurriculumRepository {
	return &CurriculumRepository{db, logger}
}

func (cr *CurriculumRepository) CreateCurriculum(ctx context.Context, tx pgx.Tx, createCurriculum *CreateCurriculum) (*Curriculum, *exceptions.AppError) {
	sqlStmt := `
		INSERT INTO curricula (class_id, subject_id, teacher_id, term_id)
		VALUES ($1, $2, $3, $4)
		RETURNING curriculum_id;
	`

	var curriculumId int32
	var err error

	if tx != nil {
		cr.logger.Debug("Creating curriculum in transaction")
		err = tx.QueryRow(ctx, sqlStmt, createCurriculum.ClassId, createCurriculum.SubjectId, createCurriculum.TeacherId, createCurriculum.TermId).
			Scan(&curriculumId)
	} else {
		cr.logger.Debug("Creating curriculum without transaction")
		err = cr.db.Pool.QueryRow(ctx, sqlStmt, createCurriculum.ClassId, createCurriculum.SubjectId, createCurriculum.TeacherId, createCurriculum.TermId).
			Scan(&curriculumId)
	}

	if err != nil {
		cr.logger.Error("Failed to create curriculum", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return cr.GetCurriculumById(ctx, tx, curriculumId)
}

func (cr *CurriculumRepository) GetCurriculumById(ctx context.Context, tx pgx.Tx, curriculumId int32) (*Curriculum, *exceptions.AppError) {
	if msg := domain_model.ValidateId(curriculumId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid curriculum ID", map[string]string{"curriculum_id": msg})
	}

	sqlStmt := `
		SELECT 
			c.curriculum_id,
			c.teacher_id,
			cl.class_id, cl.grade_level, cl.class_name,
			s.subject_id, s.subject_name,
			t.term_id, t.name
		FROM curricula c
		INNER JOIN classes cl ON c.class_id = cl.class_id
		INNER JOIN subjects s ON c.subject_id = s.subject_id
		INNER JOIN terms t ON c.term_id = t.term_id
		WHERE c.curriculum_id = $1;
	`

	var curriculum Curriculum
	var teacherId sql.NullInt32
	var class classes.Class
	var subject subjects.Subject
	var term terms.Term
	var err error

	if tx != nil {
		cr.logger.Debug("Getting curriculum by id in transaction", zap.Int32("curriculum_id", curriculumId))
		err = tx.QueryRow(ctx, sqlStmt, curriculumId).
			Scan(
				&curriculum.CurriculumId,
				&teacherId,
				&class.ClassId, &class.GradeLevel, &class.ClassName,
				&subject.SubjectId, &subject.SubjectName,
				&term.TermId, &term.Name,
			)
	} else {
		cr.logger.Debug("Getting curriculum by id without transaction", zap.Int32("curriculum_id", curriculumId))
		err = cr.db.Pool.QueryRow(ctx, sqlStmt, curriculumId).
			Scan(
				&curriculum.CurriculumId,
				&teacherId,
				&class.ClassId, &class.GradeLevel, &class.ClassName,
				&subject.SubjectId, &subject.SubjectName,
				&term.TermId, &term.Name,
			)
	}

	if err != nil {
		cr.logger.Error("Failed to get curriculum by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	curriculum.Class = &class
	curriculum.Subject = &subject
	curriculum.Term = &term

	if teacherId.Valid {
		teacherIdValue := int32(teacherId.Int32)
		curriculum.TeacherId = &teacherIdValue
	}

	return &curriculum, nil
}

func (cr *CurriculumRepository) GetCurricula(ctx context.Context, tx pgx.Tx) ([]Curriculum, *exceptions.AppError) {
	sqlStmt := `
		SELECT 
			c.curriculum_id,
			c.teacher_id,
			cl.class_id, cl.grade_level, cl.class_name,
			s.subject_id, s.subject_name,
			t.term_id, t.name
		FROM curricula c
		INNER JOIN classes cl ON c.class_id = cl.class_id
		INNER JOIN subjects s ON c.subject_id = s.subject_id
		INNER JOIN terms t ON c.term_id = t.term_id
		ORDER BY cl.grade_level, cl.class_name, s.subject_name;
	`

	var curricula []Curriculum
	var err error
	var rows pgx.Rows

	if tx != nil {
		cr.logger.Debug("Getting curricula in transaction")
		rows, err = tx.Query(ctx, sqlStmt)
	} else {
		cr.logger.Debug("Getting curricula without transaction")
		rows, err = cr.db.Pool.Query(ctx, sqlStmt)
	}

	if err != nil {
		cr.logger.Error("Failed to get curricula", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var curriculum Curriculum
		var teacherId sql.NullInt32
		var class classes.Class
		var subject subjects.Subject
		var term terms.Term

		err = rows.Scan(
			&curriculum.CurriculumId,
			&teacherId,
			&class.ClassId, &class.GradeLevel, &class.ClassName,
			&subject.SubjectId, &subject.SubjectName,
			&term.TermId, &term.Name,
		)
		if err != nil {
			cr.logger.Error("Failed to scan curriculum row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		curriculum.Class = &class
		curriculum.Subject = &subject
		curriculum.Term = &term

		if teacherId.Valid {
			teacherIdValue := int32(teacherId.Int32)
			curriculum.TeacherId = &teacherIdValue
		}

		curricula = append(curricula, curriculum)
	}

	if err = rows.Err(); err != nil {
		cr.logger.Error("Error iterating curricula rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return curricula, nil
}

func (cr *CurriculumRepository) DeleteCurriculum(ctx context.Context, tx pgx.Tx, curriculumId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(curriculumId); msg != "" {
		return exceptions.NewValidationError("Invalid curriculum ID", map[string]string{"curriculum_id": msg})
	}

	sqlStmt := "DELETE FROM curricula WHERE curriculum_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		cr.logger.Debug("Deleting curriculum in transaction")
		cmdTag, err = tx.Exec(ctx, sqlStmt, curriculumId)
	} else {
		cr.logger.Debug("Deleting curriculum without transaction")
		cmdTag, err = cr.db.Pool.Exec(ctx, sqlStmt, curriculumId)
	}
	if err != nil {
		cr.logger.Error("Failed to delete curriculum", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		cr.logger.Error("Failed to delete curriculum - curriculum not found", zap.Int32("curriculum_id", curriculumId))
		return exceptions.NewNotFoundError("Curriculum not found")
	}

	return nil
}

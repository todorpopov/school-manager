package absences

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/curricula"
	"github.com/todorpopov/school-manager/internal/domain_model/students"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/domain_model/terms"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IAbsenceRepository interface {
	CreateAbsence(ctx context.Context, tx pgx.Tx, createAbsence *CreateAbsence) (*Absence, *exceptions.AppError)
	GetAbsenceById(ctx context.Context, tx pgx.Tx, absenceId int32) (*Absence, *exceptions.AppError)
	GetAbsences(ctx context.Context, tx pgx.Tx) ([]Absence, *exceptions.AppError)
	DeleteAbsence(ctx context.Context, tx pgx.Tx, absenceId int32) *exceptions.AppError
}

type AbsenceRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewAbsenceRepository(db *persistence.Database, logger *zap.Logger) *AbsenceRepository {
	return &AbsenceRepository{db, logger}
}

func (ar *AbsenceRepository) CreateAbsence(ctx context.Context, tx pgx.Tx, createAbsence *CreateAbsence) (*Absence, *exceptions.AppError) {
	sqlStmt := `
		INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused)
		VALUES ($1, $2, $3, $4)
		RETURNING absence_id;
	`

	var absenceId int32
	var err error

	if tx != nil {
		ar.logger.Debug("Creating absence in transaction")
		err = tx.QueryRow(ctx, sqlStmt, createAbsence.StudentId, createAbsence.CurriculumId, createAbsence.AbsenceDate, createAbsence.IsExcused).
			Scan(&absenceId)
	} else {
		ar.logger.Debug("Creating absence without transaction")
		err = ar.db.Pool.QueryRow(ctx, sqlStmt, createAbsence.StudentId, createAbsence.CurriculumId, createAbsence.AbsenceDate, createAbsence.IsExcused).
			Scan(&absenceId)
	}

	if err != nil {
		ar.logger.Error("Failed to create absence", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return ar.GetAbsenceById(ctx, tx, absenceId)
}

func (ar *AbsenceRepository) GetAbsenceById(ctx context.Context, tx pgx.Tx, absenceId int32) (*Absence, *exceptions.AppError) {
	if msg := domain_model.ValidateId(absenceId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid absence ID", map[string]string{"absence_id": msg})
	}

	sqlStmt := `
		SELECT 
			a.absence_id,
			a.absence_date,
			a.is_excused,
			s.student_id, s.user_id, u.first_name, u.last_name, u.email,
			cl.class_id AS student_class_id, cl.grade_level AS student_grade_level, cl.class_name AS student_class_name,
			c.curriculum_id, c.teacher_id,
			ccl.class_id AS curriculum_class_id, ccl.grade_level AS curriculum_grade_level, ccl.class_name AS curriculum_class_name,
			subj.subject_id, subj.subject_name,
			t.term_id, t.name
		FROM absences a
		INNER JOIN students s ON a.student_id = s.student_id
		INNER JOIN users u ON s.user_id = u.user_id
		LEFT JOIN classes cl ON s.class_id = cl.class_id
		INNER JOIN curricula c ON a.curriculum_id = c.curriculum_id
		INNER JOIN classes ccl ON c.class_id = ccl.class_id
		INNER JOIN subjects subj ON c.subject_id = subj.subject_id
		INNER JOIN terms t ON c.term_id = t.term_id
		WHERE a.absence_id = $1;
	`

	var absence Absence
	var student students.Student
	var studentClassId sql.NullInt32
	var studentGradeLevel sql.NullInt32
	var studentClassName sql.NullString
	var curriculum curricula.Curriculum
	var curriculumClass classes.Class
	var subject subjects.Subject
	var term terms.Term
	var teacherId *int32
	var err error

	if tx != nil {
		ar.logger.Debug("Getting absence by id in transaction", zap.Int32("absence_id", absenceId))
		err = tx.QueryRow(ctx, sqlStmt, absenceId).
			Scan(
				&absence.AbsenceId,
				&absence.AbsenceDate,
				&absence.IsExcused,
				&student.StudentId, &student.UserId, &student.FirstName, &student.LastName, &student.Email,
				&studentClassId, &studentGradeLevel, &studentClassName,
				&curriculum.CurriculumId, &teacherId,
				&curriculumClass.ClassId, &curriculumClass.GradeLevel, &curriculumClass.ClassName,
				&subject.SubjectId, &subject.SubjectName,
				&term.TermId, &term.Name,
			)
	} else {
		ar.logger.Debug("Getting absence by id without transaction", zap.Int32("absence_id", absenceId))
		err = ar.db.Pool.QueryRow(ctx, sqlStmt, absenceId).
			Scan(
				&absence.AbsenceId,
				&absence.AbsenceDate,
				&absence.IsExcused,
				&student.StudentId, &student.UserId, &student.FirstName, &student.LastName, &student.Email,
				&studentClassId, &studentGradeLevel, &studentClassName,
				&curriculum.CurriculumId, &teacherId,
				&curriculumClass.ClassId, &curriculumClass.GradeLevel, &curriculumClass.ClassName,
				&subject.SubjectId, &subject.SubjectName,
				&term.TermId, &term.Name,
			)
	}

	if err != nil {
		ar.logger.Error("Failed to get absence by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	if studentClassId.Valid {
		student.Class = &classes.Class{
			ClassId:    studentClassId.Int32,
			GradeLevel: studentGradeLevel.Int32,
			ClassName:  studentClassName.String,
		}
	}

	curriculum.Class = &curriculumClass
	curriculum.Subject = &subject
	curriculum.Term = &term
	curriculum.TeacherId = teacherId
	absence.Student = student
	absence.Curriculum = curriculum

	return &absence, nil
}

func (ar *AbsenceRepository) GetAbsences(ctx context.Context, tx pgx.Tx) ([]Absence, *exceptions.AppError) {
	sqlStmt := `
		SELECT 
			a.absence_id,
			a.absence_date,
			a.is_excused,
			s.student_id, s.user_id, u.first_name, u.last_name, u.email,
			cl.class_id AS student_class_id, cl.grade_level AS student_grade_level, cl.class_name AS student_class_name,
			c.curriculum_id, c.teacher_id,
			ccl.class_id AS curriculum_class_id, ccl.grade_level AS curriculum_grade_level, ccl.class_name AS curriculum_class_name,
			subj.subject_id, subj.subject_name,
			t.term_id, t.name
		FROM absences a
		INNER JOIN students s ON a.student_id = s.student_id
		INNER JOIN users u ON s.user_id = u.user_id
		LEFT JOIN classes cl ON s.class_id = cl.class_id
		INNER JOIN curricula c ON a.curriculum_id = c.curriculum_id
		INNER JOIN classes ccl ON c.class_id = ccl.class_id
		INNER JOIN subjects subj ON c.subject_id = subj.subject_id
		INNER JOIN terms t ON c.term_id = t.term_id
		ORDER BY a.absence_date DESC, u.last_name, u.first_name;
	`

	var absences []Absence
	var err error
	var rows pgx.Rows

	if tx != nil {
		ar.logger.Debug("Getting absences in transaction")
		rows, err = tx.Query(ctx, sqlStmt)
	} else {
		ar.logger.Debug("Getting absences without transaction")
		rows, err = ar.db.Pool.Query(ctx, sqlStmt)
	}

	if err != nil {
		ar.logger.Error("Failed to get absences", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var absence Absence
		var student students.Student
		var studentClassId sql.NullInt32
		var studentGradeLevel sql.NullInt32
		var studentClassName sql.NullString
		var curriculum curricula.Curriculum
		var curriculumClass classes.Class
		var subject subjects.Subject
		var term terms.Term
		var teacherId *int32

		err = rows.Scan(
			&absence.AbsenceId,
			&absence.AbsenceDate,
			&absence.IsExcused,
			&student.StudentId, &student.UserId, &student.FirstName, &student.LastName, &student.Email,
			&studentClassId, &studentGradeLevel, &studentClassName,
			&curriculum.CurriculumId, &teacherId,
			&curriculumClass.ClassId, &curriculumClass.GradeLevel, &curriculumClass.ClassName,
			&subject.SubjectId, &subject.SubjectName,
			&term.TermId, &term.Name,
		)
		if err != nil {
			ar.logger.Error("Failed to scan absence row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		if studentClassId.Valid {
			student.Class = &classes.Class{
				ClassId:    studentClassId.Int32,
				GradeLevel: studentGradeLevel.Int32,
				ClassName:  studentClassName.String,
			}
		}

		curriculum.Class = &curriculumClass
		curriculum.Subject = &subject
		curriculum.Term = &term
		curriculum.TeacherId = teacherId
		absence.Student = student
		absence.Curriculum = curriculum

		absences = append(absences, absence)
	}

	if err = rows.Err(); err != nil {
		ar.logger.Error("Error iterating absences rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return absences, nil
}

func (ar *AbsenceRepository) DeleteAbsence(ctx context.Context, tx pgx.Tx, absenceId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(absenceId); msg != "" {
		return exceptions.NewValidationError("Invalid absence ID", map[string]string{"absence_id": msg})
	}

	sqlStmt := "DELETE FROM absences WHERE absence_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		ar.logger.Debug("Deleting absence in transaction")
		cmdTag, err = tx.Exec(ctx, sqlStmt, absenceId)
	} else {
		ar.logger.Debug("Deleting absence without transaction")
		cmdTag, err = ar.db.Pool.Exec(ctx, sqlStmt, absenceId)
	}
	if err != nil {
		ar.logger.Error("Failed to delete absence", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		ar.logger.Error("Failed to delete absence - absence not found", zap.Int32("absence_id", absenceId))
		return exceptions.NewNotFoundError("Absence not found")
	}

	return nil
}

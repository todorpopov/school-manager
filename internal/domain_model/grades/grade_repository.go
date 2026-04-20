package grades

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

type IGradeRepository interface {
	CreateGrade(ctx context.Context, tx pgx.Tx, createGrade *CreateGrade) (*Grade, *exceptions.AppError)
	GetGradeById(ctx context.Context, tx pgx.Tx, gradeId int32) (*Grade, *exceptions.AppError)
	GetGrades(ctx context.Context, tx pgx.Tx) ([]Grade, *exceptions.AppError)
	DeleteGrade(ctx context.Context, tx pgx.Tx, gradeId int32) *exceptions.AppError
}

type GradeRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewGradeRepository(db *persistence.Database, logger *zap.Logger) *GradeRepository {
	return &GradeRepository{db, logger}
}

func (gr *GradeRepository) CreateGrade(ctx context.Context, tx pgx.Tx, createGrade *CreateGrade) (*Grade, *exceptions.AppError) {
	sqlStmt := `
		INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date)
		VALUES ($1, $2, $3, $4)
		RETURNING grade_id;
	`

	var gradeId int32
	var err error

	if tx != nil {
		gr.logger.Debug("Creating grade in transaction")
		err = tx.QueryRow(ctx, sqlStmt, createGrade.StudentId, createGrade.CurriculumId, createGrade.GradeValue, createGrade.GradeDate).
			Scan(&gradeId)
	} else {
		gr.logger.Debug("Creating grade without transaction")
		err = gr.db.Pool.QueryRow(ctx, sqlStmt, createGrade.StudentId, createGrade.CurriculumId, createGrade.GradeValue, createGrade.GradeDate).
			Scan(&gradeId)
	}

	if err != nil {
		gr.logger.Error("Failed to create grade", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return gr.GetGradeById(ctx, tx, gradeId)
}

func (gr *GradeRepository) GetGradeById(ctx context.Context, tx pgx.Tx, gradeId int32) (*Grade, *exceptions.AppError) {
	if msg := domain_model.ValidateId(gradeId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid grade ID", map[string]string{"grade_id": msg})
	}

	sqlStmt := `
		SELECT 
			g.grade_id,
			g.grade_value,
			g.grade_date,
			s.student_id, s.user_id, u.first_name, u.last_name, u.email,
			cl.class_id AS student_class_id, cl.grade_level AS student_grade_level, cl.class_name AS student_class_name,
			c.curriculum_id, c.teacher_id,
			ccl.class_id AS curriculum_class_id, ccl.grade_level AS curriculum_grade_level, ccl.class_name AS curriculum_class_name,
			subj.subject_id, subj.subject_name,
			t.term_id, t.name
		FROM grades g
		INNER JOIN students s ON g.student_id = s.student_id
		INNER JOIN users u ON s.user_id = u.user_id
		LEFT JOIN classes cl ON s.class_id = cl.class_id
		INNER JOIN curricula c ON g.curriculum_id = c.curriculum_id
		INNER JOIN classes ccl ON c.class_id = ccl.class_id
		INNER JOIN subjects subj ON c.subject_id = subj.subject_id
		INNER JOIN terms t ON c.term_id = t.term_id
		WHERE g.grade_id = $1;
	`

	var grade Grade
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
		gr.logger.Debug("Getting grade by id in transaction", zap.Int32("grade_id", gradeId))
		err = tx.QueryRow(ctx, sqlStmt, gradeId).
			Scan(
				&grade.GradeId,
				&grade.GradeValue,
				&grade.GradeDate,
				&student.StudentId, &student.UserId, &student.FirstName, &student.LastName, &student.Email,
				&studentClassId, &studentGradeLevel, &studentClassName,
				&curriculum.CurriculumId, &teacherId,
				&curriculumClass.ClassId, &curriculumClass.GradeLevel, &curriculumClass.ClassName,
				&subject.SubjectId, &subject.SubjectName,
				&term.TermId, &term.Name,
			)
	} else {
		gr.logger.Debug("Getting grade by id without transaction", zap.Int32("grade_id", gradeId))
		err = gr.db.Pool.QueryRow(ctx, sqlStmt, gradeId).
			Scan(
				&grade.GradeId,
				&grade.GradeValue,
				&grade.GradeDate,
				&student.StudentId, &student.UserId, &student.FirstName, &student.LastName, &student.Email,
				&studentClassId, &studentGradeLevel, &studentClassName,
				&curriculum.CurriculumId, &teacherId,
				&curriculumClass.ClassId, &curriculumClass.GradeLevel, &curriculumClass.ClassName,
				&subject.SubjectId, &subject.SubjectName,
				&term.TermId, &term.Name,
			)
	}

	if err != nil {
		gr.logger.Error("Failed to get grade by id", zap.Error(err))
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
	grade.Student = &student
	grade.Curriculum = &curriculum

	return &grade, nil
}

func (gr *GradeRepository) GetGrades(ctx context.Context, tx pgx.Tx) ([]Grade, *exceptions.AppError) {
	sqlStmt := `
		SELECT 
			g.grade_id,
			g.grade_value,
			g.grade_date,
			s.student_id, s.user_id, u.first_name, u.last_name, u.email,
			cl.class_id AS student_class_id, cl.grade_level AS student_grade_level, cl.class_name AS student_class_name,
			c.curriculum_id, c.teacher_id,
			ccl.class_id AS curriculum_class_id, ccl.grade_level AS curriculum_grade_level, ccl.class_name AS curriculum_class_name,
			subj.subject_id, subj.subject_name,
			t.term_id, t.name
		FROM grades g
		INNER JOIN students s ON g.student_id = s.student_id
		INNER JOIN users u ON s.user_id = u.user_id
		LEFT JOIN classes cl ON s.class_id = cl.class_id
		INNER JOIN curricula c ON g.curriculum_id = c.curriculum_id
		INNER JOIN classes ccl ON c.class_id = ccl.class_id
		INNER JOIN subjects subj ON c.subject_id = subj.subject_id
		INNER JOIN terms t ON c.term_id = t.term_id
		ORDER BY g.grade_date DESC, u.last_name, u.first_name;
	`

	var grades []Grade
	var err error
	var rows pgx.Rows

	if tx != nil {
		gr.logger.Debug("Getting grades in transaction")
		rows, err = tx.Query(ctx, sqlStmt)
	} else {
		gr.logger.Debug("Getting grades without transaction")
		rows, err = gr.db.Pool.Query(ctx, sqlStmt)
	}

	if err != nil {
		gr.logger.Error("Failed to get grades", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var grade Grade
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
			&grade.GradeId,
			&grade.GradeValue,
			&grade.GradeDate,
			&student.StudentId, &student.UserId, &student.FirstName, &student.LastName, &student.Email,
			&studentClassId, &studentGradeLevel, &studentClassName,
			&curriculum.CurriculumId, &teacherId,
			&curriculumClass.ClassId, &curriculumClass.GradeLevel, &curriculumClass.ClassName,
			&subject.SubjectId, &subject.SubjectName,
			&term.TermId, &term.Name,
		)
		if err != nil {
			gr.logger.Error("Failed to scan grade row", zap.Error(err))
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
		grade.Student = &student
		grade.Curriculum = &curriculum

		grades = append(grades, grade)
	}

	if err = rows.Err(); err != nil {
		gr.logger.Error("Error iterating grades rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return grades, nil
}

func (gr *GradeRepository) DeleteGrade(ctx context.Context, tx pgx.Tx, gradeId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(gradeId); msg != "" {
		return exceptions.NewValidationError("Invalid grade ID", map[string]string{"grade_id": msg})
	}

	sqlStmt := "DELETE FROM grades WHERE grade_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		gr.logger.Debug("Deleting grade in transaction")
		cmdTag, err = tx.Exec(ctx, sqlStmt, gradeId)
	} else {
		gr.logger.Debug("Deleting grade without transaction")
		cmdTag, err = gr.db.Pool.Exec(ctx, sqlStmt, gradeId)
	}
	if err != nil {
		gr.logger.Error("Failed to delete grade", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		gr.logger.Error("Failed to delete grade - grade not found", zap.Int32("grade_id", gradeId))
		return exceptions.NewNotFoundError("Grade not found")
	}

	return nil
}

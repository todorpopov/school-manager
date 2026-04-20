package students

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/parents"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

const MaxParentsPerStudent = 2

type IStudentParentRepository interface {
	LinkParentToStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError
	UnlinkParentFromStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError
	GetParentsForStudent(ctx context.Context, tx pgx.Tx, studentId int32) ([]parents.Parent, *exceptions.AppError)
	GetStudentsForParent(ctx context.Context, tx pgx.Tx, parentId int32) ([]Student, *exceptions.AppError)
	CountParentsForStudent(ctx context.Context, tx pgx.Tx, studentId int32) (int, *exceptions.AppError)
}

type StudentParentRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewStudentParentRepository(db *persistence.Database, logger *zap.Logger) *StudentParentRepository {
	return &StudentParentRepository{db, logger}
}

func (spr *StudentParentRepository) LinkParentToStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	sqlStmt := `
		INSERT INTO student_parents (parent_id, student_id)
		VALUES ($1, $2);
	`

	var err error
	if tx != nil {
		spr.logger.Debug("Linking parent to student in transaction", zap.Int32("student_id", studentId), zap.Int32("parent_id", parentId))
		_, err = tx.Exec(ctx, sqlStmt, parentId, studentId)
	} else {
		spr.logger.Debug("Linking parent to student without transaction", zap.Int32("student_id", studentId), zap.Int32("parent_id", parentId))
		_, err = spr.db.Pool.Exec(ctx, sqlStmt, parentId, studentId)
	}

	if err != nil {
		spr.logger.Error("Failed to link parent to student", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	return nil
}

func (spr *StudentParentRepository) UnlinkParentFromStudent(ctx context.Context, tx pgx.Tx, studentId int32, parentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	sqlStmt := `
		DELETE FROM student_parents
		WHERE parent_id = $1 AND student_id = $2;
	`

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		spr.logger.Debug("Unlinking parent from student in transaction", zap.Int32("student_id", studentId), zap.Int32("parent_id", parentId))
		cmdTag, err = tx.Exec(ctx, sqlStmt, parentId, studentId)
	} else {
		spr.logger.Debug("Unlinking parent from student without transaction", zap.Int32("student_id", studentId), zap.Int32("parent_id", parentId))
		cmdTag, err = spr.db.Pool.Exec(ctx, sqlStmt, parentId, studentId)
	}

	if err != nil {
		spr.logger.Error("Failed to unlink parent from student", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		spr.logger.Error("Failed to unlink parent from student - relationship not found")
		return exceptions.NewNotFoundError("Parent-student relationship not found")
	}

	return nil
}

func (spr *StudentParentRepository) GetParentsForStudent(ctx context.Context, tx pgx.Tx, studentId int32) ([]parents.Parent, *exceptions.AppError) {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	sqlStmt := `
		SELECT p.parent_id, p.user_id, u.first_name, u.last_name, u.email
		FROM student_parents sp
		INNER JOIN parents p ON sp.parent_id = p.parent_id
		INNER JOIN users u ON p.user_id = u.user_id
		WHERE sp.student_id = $1
		ORDER BY p.parent_id;
	`

	var studentParents []parents.Parent
	var err error
	var rows pgx.Rows

	if tx != nil {
		spr.logger.Debug("Getting parents for student in transaction", zap.Int32("student_id", studentId))
		rows, err = tx.Query(ctx, sqlStmt, studentId)
	} else {
		spr.logger.Debug("Getting parents for student without transaction", zap.Int32("student_id", studentId))
		rows, err = spr.db.Pool.Query(ctx, sqlStmt, studentId)
	}

	if err != nil {
		spr.logger.Error("Failed to get parents for student", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var parent parents.Parent
		err = rows.Scan(&parent.ParentId, &parent.UserId, &parent.FirstName, &parent.LastName, &parent.Email)
		if err != nil {
			spr.logger.Error("Failed to scan parent row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}
		studentParents = append(studentParents, parent)
	}

	if err = rows.Err(); err != nil {
		spr.logger.Error("Error iterating parent rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return studentParents, nil
}

func (spr *StudentParentRepository) GetStudentsForParent(ctx context.Context, tx pgx.Tx, parentId int32) ([]Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(parentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid parent ID", map[string]string{"parent_id": msg})
	}

	sqlStmt := `
		SELECT 
			s.student_id, s.user_id,
			u.first_name, u.last_name, u.email,
			c.class_id, c.grade_level, c.class_name
		FROM student_parents sp
		INNER JOIN students s ON sp.student_id = s.student_id
		INNER JOIN users u ON s.user_id = u.user_id
		LEFT JOIN classes c ON s.class_id = c.class_id
		WHERE sp.parent_id = $1
		ORDER BY s.student_id;
	`

	var parentStudents []Student
	var err error
	var rows pgx.Rows

	if tx != nil {
		spr.logger.Debug("Getting students for parent in transaction", zap.Int32("parent_id", parentId))
		rows, err = tx.Query(ctx, sqlStmt, parentId)
	} else {
		spr.logger.Debug("Getting students for parent without transaction", zap.Int32("parent_id", parentId))
		rows, err = spr.db.Pool.Query(ctx, sqlStmt, parentId)
	}

	if err != nil {
		spr.logger.Error("Failed to get students for parent", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		var classId *int32
		var gradeLevel *int32
		var className *string

		err = rows.Scan(
			&student.StudentId, &student.UserId,
			&student.FirstName, &student.LastName, &student.Email,
			&classId, &gradeLevel, &className,
		)
		if err != nil {
			spr.logger.Error("Failed to scan student row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		if classId != nil && gradeLevel != nil && className != nil {
			student.Class = &classes.Class{
				ClassId:    *classId,
				GradeLevel: *gradeLevel,
				ClassName:  *className,
			}
		}

		parentStudents = append(parentStudents, student)
	}

	if err = rows.Err(); err != nil {
		spr.logger.Error("Error iterating student rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return parentStudents, nil
}

func (spr *StudentParentRepository) CountParentsForStudent(ctx context.Context, tx pgx.Tx, studentId int32) (int, *exceptions.AppError) {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return 0, exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	sqlStmt := `
		SELECT COUNT(*)
		FROM student_parents
		WHERE student_id = $1;
	`

	var count int
	var err error

	if tx != nil {
		spr.logger.Debug("Counting parents for student in transaction", zap.Int32("student_id", studentId))
		err = tx.QueryRow(ctx, sqlStmt, studentId).Scan(&count)
	} else {
		spr.logger.Debug("Counting parents for student without transaction", zap.Int32("student_id", studentId))
		err = spr.db.Pool.QueryRow(ctx, sqlStmt, studentId).Scan(&count)
	}

	if err != nil {
		spr.logger.Error("Failed to count parents for student", zap.Error(err))
		return 0, exceptions.PgErrorToAppError(err)
	}

	return count, nil
}



package students

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/persistence"
	"go.uber.org/zap"
)

type IStudentRepository interface {
	CreateStudent(ctx context.Context, tx pgx.Tx, schoolId int32, userId int32, classId *int32) (*Student, *exceptions.AppError)
	GetStudentById(ctx context.Context, tx pgx.Tx, studentId int32) (*Student, *exceptions.AppError)
	GetStudentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Student, *exceptions.AppError)
	GetStudents(ctx context.Context, tx pgx.Tx) ([]Student, *exceptions.AppError)
	UpdateStudentClass(ctx context.Context, tx pgx.Tx, studentId int32, classId *int32) *exceptions.AppError
	DeleteStudent(ctx context.Context, tx pgx.Tx, studentId int32) *exceptions.AppError
}

type StudentRepository struct {
	db     *persistence.Database
	logger *zap.Logger
}

func NewStudentRepository(db *persistence.Database, logger *zap.Logger) *StudentRepository {
	return &StudentRepository{db, logger}
}

func (sr *StudentRepository) CreateStudent(ctx context.Context, tx pgx.Tx, schoolId int32, userId int32, classId *int32) (*Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(schoolId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid school ID", map[string]string{"school_id": msg})
	}
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}
	if classId != nil {
		if msg := domain_model.ValidateId(*classId); msg != "" {
			return nil, exceptions.NewValidationError("Invalid class ID", map[string]string{"class_id": msg})
		}
	}

	sql := `
		INSERT INTO students (school_id, user_id, class_id)
		VALUES ($1, $2, $3)
		RETURNING student_id;
	`

	var studentId int32
	var err error

	if tx != nil {
		sr.logger.Debug("Creating student in transaction")
		err = tx.QueryRow(ctx, sql, schoolId, userId, classId).Scan(&studentId)
	} else {
		sr.logger.Debug("Creating student without transaction")
		err = sr.db.Pool.QueryRow(ctx, sql, schoolId, userId, classId).Scan(&studentId)
	}

	if err != nil {
		sr.logger.Error("Failed to create student", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return sr.GetStudentById(ctx, tx, studentId)
}

func (sr *StudentRepository) GetStudentById(ctx context.Context, tx pgx.Tx, studentId int32) (*Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	sql := `
		SELECT 
			s.student_id, s.user_id, s.class_id,
			u.first_name, u.last_name, u.email,
			sch.school_id, sch.school_name, sch.school_address,
			c.class_id, c.grade_level, c.class_name
		FROM students s
		INNER JOIN users u ON s.user_id = u.user_id
		INNER JOIN schools sch ON s.school_id = sch.school_id
		LEFT JOIN classes c ON s.class_id = c.class_id
		WHERE s.student_id = $1;
	`

	var student Student
	var class classes.Class
	var school schools.School
	var tempClassId *int32
	var classId *int32
	var gradeLevel *int32
	var className *string
	var err error

	if tx != nil {
		sr.logger.Debug("Getting student by id in transaction", zap.Int32("student_id", studentId))
		err = tx.QueryRow(ctx, sql, studentId).
			Scan(
				&student.StudentId, &student.UserId, &tempClassId,
				&student.FirstName, &student.LastName, &student.Email,
				&school.SchoolId, &school.SchoolName, &school.SchoolAddress,
				&classId, &gradeLevel, &className,
			)
	} else {
		sr.logger.Debug("Getting student by id without transaction", zap.Int32("student_id", studentId))
		err = sr.db.Pool.QueryRow(ctx, sql, studentId).
			Scan(
				&student.StudentId, &student.UserId, &tempClassId,
				&student.FirstName, &student.LastName, &student.Email,
				&school.SchoolId, &school.SchoolName, &school.SchoolAddress,
				&classId, &gradeLevel, &className,
			)
	}

	if err != nil {
		sr.logger.Error("Failed to get student by id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	student.School = &school

	if classId != nil && gradeLevel != nil && className != nil {
		class.ClassId = *classId
		class.GradeLevel = *gradeLevel
		class.ClassName = *className
		student.Class = &class
	}

	return &student, nil
}

func (sr *StudentRepository) GetStudentByUserId(ctx context.Context, tx pgx.Tx, userId int32) (*Student, *exceptions.AppError) {
	if msg := domain_model.ValidateId(userId); msg != "" {
		return nil, exceptions.NewValidationError("Invalid user ID", map[string]string{"user_id": msg})
	}

	sql := `
		SELECT 
			s.student_id, s.user_id, s.class_id,
			u.first_name, u.last_name, u.email,
			sch.school_id, sch.school_name, sch.school_address,
			c.class_id, c.grade_level, c.class_name
		FROM students s
		INNER JOIN users u ON s.user_id = u.user_id
		INNER JOIN schools sch ON s.school_id = sch.school_id
		LEFT JOIN classes c ON s.class_id = c.class_id
		WHERE s.user_id = $1;
	`

	var student Student
	var class classes.Class
	var school schools.School
	var tempClassId *int32
	var classId *int32
	var gradeLevel *int32
	var className *string
	var err error

	if tx != nil {
		sr.logger.Debug("Getting student by user_id in transaction", zap.Int32("user_id", userId))
		err = tx.QueryRow(ctx, sql, userId).
			Scan(
				&student.StudentId, &student.UserId, &tempClassId,
				&student.FirstName, &student.LastName, &student.Email,
				&school.SchoolId, &school.SchoolName, &school.SchoolAddress,
				&classId, &gradeLevel, &className,
			)
	} else {
		sr.logger.Debug("Getting student by user_id without transaction", zap.Int32("user_id", userId))
		err = sr.db.Pool.QueryRow(ctx, sql, userId).
			Scan(
				&student.StudentId, &student.UserId, &tempClassId,
				&student.FirstName, &student.LastName, &student.Email,
				&school.SchoolId, &school.SchoolName, &school.SchoolAddress,
				&classId, &gradeLevel, &className,
			)
	}

	if err != nil {
		sr.logger.Error("Failed to get student by user_id", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	student.School = &school

	if classId != nil && gradeLevel != nil && className != nil {
		class.ClassId = *classId
		class.GradeLevel = *gradeLevel
		class.ClassName = *className
		student.Class = &class
	}

	return &student, nil
}

func (sr *StudentRepository) GetStudents(ctx context.Context, tx pgx.Tx) ([]Student, *exceptions.AppError) {
	sql := `
		SELECT 
			s.student_id, s.user_id, s.class_id,
			u.first_name, u.last_name, u.email,
			sch.school_id, sch.school_name, sch.school_address,
			c.class_id, c.grade_level, c.class_name
		FROM students s
		INNER JOIN users u ON s.user_id = u.user_id
		INNER JOIN schools sch ON s.school_id = sch.school_id
		LEFT JOIN classes c ON s.class_id = c.class_id
		ORDER BY c.grade_level, c.class_name, u.last_name, u.first_name;
	`

	var students []Student
	var err error
	var rows pgx.Rows

	if tx != nil {
		sr.logger.Debug("Getting students in transaction")
		rows, err = tx.Query(ctx, sql)
	} else {
		sr.logger.Debug("Getting students without transaction")
		rows, err = sr.db.Pool.Query(ctx, sql)
	}

	if err != nil {
		sr.logger.Error("Failed to get students", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		var class classes.Class
		var school schools.School
		var tempClassId *int32
		var classId *int32
		var gradeLevel *int32
		var className *string

		err = rows.Scan(
			&student.StudentId, &student.UserId, &tempClassId,
			&student.FirstName, &student.LastName, &student.Email,
			&school.SchoolId, &school.SchoolName, &school.SchoolAddress,
			&classId, &gradeLevel, &className,
		)
		if err != nil {
			sr.logger.Error("Failed to scan student row", zap.Error(err))
			return nil, exceptions.PgErrorToAppError(err)
		}

		student.School = &school

		if classId != nil && gradeLevel != nil && className != nil {
			class.ClassId = *classId
			class.GradeLevel = *gradeLevel
			class.ClassName = *className
			student.Class = &class
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		sr.logger.Error("Error iterating students rows", zap.Error(err))
		return nil, exceptions.PgErrorToAppError(err)
	}

	return students, nil
}

func (sr *StudentRepository) UpdateStudentClass(ctx context.Context, tx pgx.Tx, studentId int32, classId *int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}
	if classId != nil {
		if msg := domain_model.ValidateId(*classId); msg != "" {
			return exceptions.NewValidationError("Invalid class ID", map[string]string{"class_id": msg})
		}
	}

	sql := `
		UPDATE students
		SET class_id = $1
		WHERE student_id = $2;
	`

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		sr.logger.Debug("Updating student class in transaction", zap.Int32("student_id", studentId))
		cmdTag, err = tx.Exec(ctx, sql, classId, studentId)
	} else {
		sr.logger.Debug("Updating student class without transaction", zap.Int32("student_id", studentId))
		cmdTag, err = sr.db.Pool.Exec(ctx, sql, classId, studentId)
	}

	if err != nil {
		sr.logger.Error("Failed to update student class", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		sr.logger.Error("Failed to update student class - student not found", zap.Int32("student_id", studentId))
		return exceptions.NewNotFoundError("Student not found")
	}

	return nil
}

func (sr *StudentRepository) DeleteStudent(ctx context.Context, tx pgx.Tx, studentId int32) *exceptions.AppError {
	if msg := domain_model.ValidateId(studentId); msg != "" {
		return exceptions.NewValidationError("Invalid student ID", map[string]string{"student_id": msg})
	}

	sql := "DELETE FROM students WHERE student_id = $1;"

	var err error
	var cmdTag pgconn.CommandTag
	if tx != nil {
		sr.logger.Debug("Deleting student in transaction")
		cmdTag, err = tx.Exec(ctx, sql, studentId)
	} else {
		sr.logger.Debug("Deleting student without transaction")
		cmdTag, err = sr.db.Pool.Exec(ctx, sql, studentId)
	}
	if err != nil {
		sr.logger.Error("Failed to delete student", zap.Error(err))
		return exceptions.PgErrorToAppError(err)
	}

	if cmdTag.RowsAffected() == 0 {
		sr.logger.Error("Failed to delete student - student not found", zap.Int32("student_id", studentId))
		return exceptions.NewNotFoundError("Student not found")
	}

	return nil
}

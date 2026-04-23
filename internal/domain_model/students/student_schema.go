package students

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Student struct {
	StudentId int32           `json:"student_id"`
	UserId    int32           `json:"user_id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     string          `json:"email"`
	School    *schools.School `json:"school"`
	Class     *classes.Class  `json:"class,omitempty"`
	Roles     []string        `json:"roles,omitempty"`
}

type CreateStudent struct {
	SchoolId  int32  `json:"school_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	ClassId   *int32 `json:"class_id"`
}

type UpdateStudent struct {
	StudentId int32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ClassId   *int32 `json:"class_id"`
}

type UpdateStudentRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ClassId   *int32 `json:"class_id"`
}

func ValidateCreateStudent(createStudent *CreateStudent) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(createStudent.SchoolId)
	if msg != "" {
		messages["school_id"] = msg
	}

	msg = domain_model.ValidateString(&createStudent.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&createStudent.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&createStudent.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	msg = domain_model.ValidatePassword(&createStudent.Password, true)
	if msg != "" {
		messages["password"] = msg
	}

	if createStudent.ClassId != nil {
		msg = domain_model.ValidateId(*createStudent.ClassId)
		if msg != "" {
			messages["class_id"] = msg
		}
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during student creation", messages)
	}
	return nil
}

func ValidateUpdateStudent(updateStudent *UpdateStudent) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateStudent.StudentId)
	if msg != "" {
		messages["student_id"] = msg
	}

	msg = domain_model.ValidateString(&updateStudent.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&updateStudent.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&updateStudent.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	if updateStudent.ClassId != nil {
		msg = domain_model.ValidateId(*updateStudent.ClassId)
		if msg != "" {
			messages["class_id"] = msg
		}
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during student update", messages)
	}
	return nil
}

package teachers

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Teacher struct {
	TeacherId int32    `json:"teacher_id"`
	UserId    int32    `json:"user_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles,omitempty"`
}

type CreateTeacher struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateTeacher struct {
	TeacherId int32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateTeacherRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func ValidateCreateTeacher(createTeacher *CreateTeacher) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createTeacher.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&createTeacher.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&createTeacher.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	msg = domain_model.ValidatePassword(&createTeacher.Password, true)
	if msg != "" {
		messages["password"] = msg
	}


	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during teacher creation", messages)
	}
	return nil
}

func ValidateUpdateTeacher(updateTeacher *UpdateTeacher) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateTeacher.TeacherId)
	if msg != "" {
		messages["teacher_id"] = msg
	}

	msg = domain_model.ValidateString(&updateTeacher.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&updateTeacher.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&updateTeacher.Email, true)
	if msg != "" {
		messages["email"] = msg
	}


	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during teacher update", messages)
	}
	return nil
}

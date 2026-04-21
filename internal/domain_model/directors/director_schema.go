package directors

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Director struct {
	DirectorId int32    `json:"director_id"`
	UserId     int32    `json:"user_id"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Roles      []string `json:"roles,omitempty"`
}

type CreateDirector struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateDirector struct {
	DirectorId int32
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
}

type UpdateDirectorRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func ValidateCreateDirector(createDirector *CreateDirector) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createDirector.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&createDirector.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&createDirector.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	msg = domain_model.ValidatePassword(&createDirector.Password, true)
	if msg != "" {
		messages["password"] = msg
	}


	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during director creation", messages)
	}
	return nil
}

func ValidateUpdateDirector(updateDirector *UpdateDirector) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateDirector.DirectorId)
	if msg != "" {
		messages["director_id"] = msg
	}

	msg = domain_model.ValidateString(&updateDirector.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&updateDirector.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&updateDirector.Email, true)
	if msg != "" {
		messages["email"] = msg
	}


	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during director update", messages)
	}
	return nil
}


package parents

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Parent struct {
	ParentId  int32    `json:"parent_id"`
	UserId    int32    `json:"user_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Roles     []string `json:"roles,omitempty"`
}

type CreateParent struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateParent struct {
	ParentId  int32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateParentRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func ValidateCreateParent(createParent *CreateParent) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createParent.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&createParent.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&createParent.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	msg = domain_model.ValidatePassword(&createParent.Password, true)
	if msg != "" {
		messages["password"] = msg
	}


	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during parent creation", messages)
	}
	return nil
}

func ValidateUpdateParent(updateParent *UpdateParent) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateParent.ParentId)
	if msg != "" {
		messages["parent_id"] = msg
	}

	msg = domain_model.ValidateString(&updateParent.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&updateParent.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&updateParent.Email, true)
	if msg != "" {
		messages["email"] = msg
	}


	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during parent update", messages)
	}
	return nil
}

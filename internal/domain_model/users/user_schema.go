package users

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type User struct {
	UserId    int32   `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Password  *string `json:"password,omitempty"`
}

type CreateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUser struct {
	UserId    int32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserPassword struct {
	UserId   int32  `json:"user_id"`
	Password string `json:"password"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password"`
}

func ValidateCreateUser(createUser *CreateUser) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createUser.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&createUser.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&createUser.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	msg = domain_model.ValidatePassword(&createUser.Password, true)
	if msg != "" {
		messages["password"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during user creation", messages)
	}
	return nil
}

func ValidateUpdateUser(updateUser *UpdateUser) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateUser.UserId)
	if msg != "" {
		messages["user_id"] = msg
	}

	msg = domain_model.ValidateString(&updateUser.FirstName, 1, 255, true)
	if msg != "" {
		messages["first_name"] = msg
	}

	msg = domain_model.ValidateString(&updateUser.LastName, 1, 255, true)
	if msg != "" {
		messages["last_name"] = msg
	}

	msg = domain_model.ValidateEmail(&updateUser.Email, true)
	if msg != "" {
		messages["email"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during user update", messages)
	}
	return nil
}

func ValidateUpdateUserPassword(updateUserPass *UpdateUserPassword) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateUserPass.UserId)
	if msg != "" {
		messages["user_id"] = msg
	}

	msg = domain_model.ValidatePassword(&updateUserPass.Password, true)
	if msg != "" {
		messages["password"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during user password update", messages)
	}
	return nil
}

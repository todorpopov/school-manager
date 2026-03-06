package roles

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Role struct {
	RoleId   int32  `json:"role_id"`
	RoleName string `json:"role_name"`
}

type CreateRole struct {
	RoleName string `json:"role_name"`
}

func ValidateCreateRole(createRole *CreateRole) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateRoleName(createRole.RoleName)
	if msg != "" {
		messages["role_name"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed for creating role", messages)
	}
	return nil
}

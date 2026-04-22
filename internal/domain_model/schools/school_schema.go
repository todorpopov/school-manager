package schools

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type School struct {
	SchoolId      int32  `json:"school_id"`
	SchoolName    string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
}

type CreateSchool struct {
	SchoolName    string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
}

type UpdateSchool struct {
	SchoolId      int32
	SchoolName    string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
}

type UpdateSchoolRequest struct {
	SchoolName    string `json:"school_name"`
	SchoolAddress string `json:"school_address"`
}

func ValidateCreateSchool(createSchool *CreateSchool) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createSchool.SchoolName, 1, 255, true)
	if msg != "" {
		messages["school_name"] = msg
	}

	msg = domain_model.ValidateString(&createSchool.SchoolAddress, 1, 255, true)
	if msg != "" {
		messages["school_address"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during school creation", messages)
	}
	return nil
}

func ValidateUpdateSchool(updateSchool *UpdateSchool) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(updateSchool.SchoolId)
	if msg != "" {
		messages["school_id"] = msg
	}

	msg = domain_model.ValidateString(&updateSchool.SchoolName, 1, 255, true)
	if msg != "" {
		messages["school_name"] = msg
	}

	msg = domain_model.ValidateString(&updateSchool.SchoolAddress, 1, 255, true)
	if msg != "" {
		messages["school_address"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during school update", messages)
	}
	return nil
}


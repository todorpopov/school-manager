package classes

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Class struct {
	ClassId    int32  `json:"class_id"`
	SchoolId   int32  `json:"school_id"`
	GradeLevel int32  `json:"grade_level"`
	ClassName  string `json:"class_name"`
}

type CreateClass struct {
	SchoolId   int32  `json:"school_id"`
	GradeLevel int32  `json:"grade_level"`
	ClassName  string `json:"class_name"`
}

func ValidateCreateClass(createClass *CreateClass) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(createClass.SchoolId)
	if msg != "" {
		messages["school_id"] = msg
	}

	if createClass.GradeLevel < 1 || createClass.GradeLevel > 12 {
		messages["grade_level"] = "Grade level must be between 1 and 12"
	}

	msg = domain_model.ValidateString(&createClass.ClassName, 1, 255, true)
	if msg != "" {
		messages["class_name"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during class creation", messages)
	}
	return nil
}

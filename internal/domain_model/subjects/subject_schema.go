package subjects

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Subject struct {
	SubjectId   int32  `json:"subject_id"`
	SubjectName string `json:"subject_name"`
}

type CreateSubject struct {
	SubjectName string `json:"subject_name"`
}

func ValidateCreateSubject(createSubject *CreateSubject) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createSubject.SubjectName, 1, 255, true)
	if msg != "" {
		messages["subject_name"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during subject creation", messages)
	}
	return nil
}


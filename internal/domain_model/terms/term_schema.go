package terms

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Term struct {
	TermId int32  `json:"term_id"`
	Name   string `json:"name"`
}

type CreateTerm struct {
	Name string `json:"name"`
}

func ValidateCreateTerm(createTerm *CreateTerm) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createTerm.Name, 1, 255, true)
	if msg != "" {
		messages["name"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during term creation", messages)
	}
	return nil
}


package terms

import (
	"time"

	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Term struct {
	TermId    int32     `json:"term_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type CreateTerm struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func ValidateCreateTerm(createTerm *CreateTerm) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateString(&createTerm.Name, 1, 255, true)
	if msg != "" {
		messages["name"] = msg
	}

	msg = domain_model.ValidateString(&createTerm.StartDate, 1, 10, true)
	if msg != "" {
		messages["start_date"] = msg
	}

	msg = domain_model.ValidateIsoDate(&createTerm.StartDate, true)
	if msg != "" {
		messages["start_date"] = msg
	}

	msg = domain_model.ValidateString(&createTerm.EndDate, 1, 10, true)
	if msg != "" {
		messages["end_date"] = msg
	}

	msg = domain_model.ValidateIsoDate(&createTerm.EndDate, true)
	if msg != "" {
		messages["end_date"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during term creation", messages)
	}
	return nil
}


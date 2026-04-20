package absences

import (
	"time"

	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/curricula"
	"github.com/todorpopov/school-manager/internal/domain_model/students"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Absence struct {
	AbsenceId   int32                `json:"absence_id"`
	Student     students.Student     `json:"student"`
	Curriculum  curricula.Curriculum `json:"curriculum"`
	AbsenceDate time.Time            `json:"absence_date"`
	IsExcused   bool                 `json:"is_excused"`
}

type CreateAbsence struct {
	StudentId    int32  `json:"student_id"`
	CurriculumId int32  `json:"curriculum_id"`
	AbsenceDate  string `json:"absence_date"`
	IsExcused    bool   `json:"is_excused"`
}

func ValidateCreateAbsence(createAbsence *CreateAbsence) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(createAbsence.StudentId)
	if msg != "" {
		messages["student_id"] = msg
	}

	msg = domain_model.ValidateId(createAbsence.CurriculumId)
	if msg != "" {
		messages["curriculum_id"] = msg
	}

	msg = domain_model.ValidateString(&createAbsence.AbsenceDate, 1, 10, true)
	if msg != "" {
		messages["absence_date"] = msg
	}

	msg = domain_model.ValidateIsoDate(&createAbsence.AbsenceDate, true)
	if msg != "" {
		messages["absence_date"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during absence creation", messages)
	}
	return nil
}


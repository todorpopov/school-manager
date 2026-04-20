package grades

import (
	"time"

	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/curricula"
	"github.com/todorpopov/school-manager/internal/domain_model/students"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Grade struct {
	GradeId    int32                 `json:"grade_id"`
	Student    *students.Student     `json:"student"`
	Curriculum *curricula.Curriculum `json:"curriculum"`
	GradeValue float32               `json:"grade_value"`
	GradeDate  time.Time             `json:"grade_date"`
}

type CreateGrade struct {
	StudentId    int32   `json:"student_id"`
	CurriculumId int32   `json:"curriculum_id"`
	GradeValue   float32 `json:"grade_value"`
	GradeDate    string  `json:"grade_date"`
}

func ValidateCreateGrade(createGrade *CreateGrade) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(createGrade.StudentId)
	if msg != "" {
		messages["student_id"] = msg
	}

	msg = domain_model.ValidateId(createGrade.CurriculumId)
	if msg != "" {
		messages["curriculum_id"] = msg
	}

	if createGrade.GradeValue < 2.00 || createGrade.GradeValue > 6.00 {
		messages["grade_value"] = "Grade value must be between 2.00 and 6.00"
	}

	msg = domain_model.ValidateString(&createGrade.GradeDate, 1, 10, true)
	if msg != "" {
		messages["grade_date"] = msg
	}

	msg = domain_model.ValidateIsoDate(&createGrade.GradeDate, true)
	if msg != "" {
		messages["grade_date"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during grade creation", messages)
	}
	return nil
}

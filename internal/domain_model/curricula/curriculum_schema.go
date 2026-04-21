package curricula

import (
	"github.com/todorpopov/school-manager/internal/domain_model"
	"github.com/todorpopov/school-manager/internal/domain_model/classes"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/domain_model/terms"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Curriculum struct {
	CurriculumId int32              `json:"curriculum_id"`
	Class        *classes.Class     `json:"class"`
	Subject      *subjects.Subject  `json:"subject"`
	TeacherId    *int32             `json:"teacher_id,omitempty"`
	Term         *terms.Term        `json:"term"`
}

type CreateCurriculum struct {
	ClassId   int32 `json:"class_id"`
	SubjectId int32 `json:"subject_id"`
	TeacherId int32 `json:"teacher_id"`
	TermId    int32 `json:"term_id"`
}

func ValidateCreateCurriculum(createCurriculum *CreateCurriculum) *exceptions.AppError {
	messages := map[string]string{}
	var msg string

	msg = domain_model.ValidateId(createCurriculum.ClassId)
	if msg != "" {
		messages["class_id"] = msg
	}

	msg = domain_model.ValidateId(createCurriculum.SubjectId)
	if msg != "" {
		messages["subject_id"] = msg
	}

	msg = domain_model.ValidateId(createCurriculum.TeacherId)
	if msg != "" {
		messages["teacher_id"] = msg
	}

	msg = domain_model.ValidateId(createCurriculum.TermId)
	if msg != "" {
		messages["term_id"] = msg
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Validation failed during curriculum creation", messages)
	}
	return nil
}


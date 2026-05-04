package reporting

import (
	"github.com/todorpopov/school-manager/internal/domain_model/schools"
	"github.com/todorpopov/school-manager/internal/domain_model/subjects"
	"github.com/todorpopov/school-manager/internal/domain_model/teachers"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

const (
	ReportTypeGrades   string = "grades"
	ReportTypeAbsences string = "absences"
)

type DynamicReportingResponse interface{}

type AbsencesReport []SchoolAbsences

type SchoolAbsences struct {
	School   schools.School            `json:"school"`
	Teachers []TeacherSubjectsAbsences `json:"teachers"`
}

type TeacherSubjectsAbsences struct {
	Teacher  teachers.Teacher        `json:"teacher"`
	Subjects []SubjectAbsenceSummary `json:"subjects,omitempty"`
}

type SubjectAbsenceSummary struct {
	Subject        subjects.Subject `json:"subject"`
	AbsenceSummary AbsenceSummary   `json:"absence_summary"`
}

type AbsenceSummary struct {
	NumberOfAbsences   int32 `json:"number_of_absences"`
	NumberOfExcuses    int32 `json:"number_of_excuses"`
	NumberOfNonExcuses int32 `json:"number_of_non_excuses"`
}

type GradesReport []SchoolGrades

type SchoolGrades struct {
	School   schools.School          `json:"school"`
	Teachers []TeacherSubjectsGrades `json:"teachers"`
}

type TeacherSubjectsGrades struct {
	Teacher  teachers.Teacher      `json:"teacher"`
	Subjects []SubjectGradeSummary `json:"subjects,omitempty"`
}

type SubjectGradeSummary struct {
	Subject      subjects.Subject `json:"subject"`
	GradeSummary GradeSummary     `json:"absence_summary"`
}

type GradeSummary struct {
	NumberOfGrades    int32             `json:"number_of_grades"`
	GradeDistribution GradeDistribution `json:"grade_distribution"`
}

type GradeDistribution struct {
	NumberOfSixes  int32 `json:"number_of_sixes"`
	NumberOfFives  int32 `json:"number_of_fives"`
	NumberOfFours  int32 `json:"number_of_fours"`
	NumberOfThrees int32 `json:"number_of_threes"`
	NumberOfTwos   int32 `json:"number_of_twos"`
}

type ReportQueryRequest struct {
	ReportType string  `json:"report_type"`
	SchoolIds  []int32 `json:"school_ids"`
	TeacherIds []int32 `json:"teacher_ids"`
	SubjectIds []int32 `json:"subject_ids"`
}

func ValidateReportQueryRequest(reportQueryRequest *ReportQueryRequest) *exceptions.AppError {
	messages := map[string]string{}

	if reportQueryRequest.ReportType != ReportTypeGrades && reportQueryRequest.ReportType != ReportTypeAbsences {
		messages["report_type"] = "Invalid report type! Please choose either 'grades' or 'absences'."
	}

	if len(messages) > 0 {
		return exceptions.NewValidationError("Report query validation error", messages)
	}
	return nil
}

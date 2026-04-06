package curricula

type Curriculum struct {
	CurriculumId int32 `json:"curriculum_id"`
	ClassId      int32 `json:"class_id"`
	SubjectId    int32 `json:"subject_id"`
	TeacherId    int32 `json:"teacher_id"`
	TermId       int32 `json:"term_id"`
}

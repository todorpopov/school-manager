package grades

type Grade struct {
	GradeId      int32  `json:"grade_id"`
	StudentId    int32  `json:"student_id"`
	CurriculumId int32  `json:"curriculum_id"`
	TeacherId    int32  `json:"teacher_id"`
	GradeValue   int32  `json:"grade_value"`
	GradeDate    string `json:"grade_date"`
}

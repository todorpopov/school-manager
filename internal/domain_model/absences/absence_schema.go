package absences

type Absence struct {
	AbsenceId    int32  `json:"absence_id"`
	StudentId    int32  `json:"student_id"`
	CurriculumId int32  `json:"curriculum_id"`
	TeacherId    int32  `json:"teacher_id"`
	AbsenceDate  string `json:"absence_date"`
	IsExcused    bool   `json:"is_excused"`
}

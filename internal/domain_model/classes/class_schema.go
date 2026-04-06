package classes

type Class struct {
	ClassId    int32  `json:"class_id"`
	GradeLevel int32  `json:"grade_level"`
	ClassName  string `json:"class_name"`
}

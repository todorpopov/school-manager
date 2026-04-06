package students

type Student struct {
	StudentId int32 `json:"student_id"`
	UserId    int32 `json:"user_id"`
	ClassId   int32 `json:"class_id"`
}

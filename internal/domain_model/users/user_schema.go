package users

type User struct {
	UserId    int32   `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Password  *string `json:"password,omitempty"`
}

type CreateUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUser struct {
	UserId    int32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UpdateUserPassword struct {
	UserId   int32  `json:"user_id"`
	Password string `json:"password"`
}

type UpdateUserPasswordRequest struct {
	Password string `json:"password"`
}

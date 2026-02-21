package internal

type ApiResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApiResponse(error bool, message string, data interface{}) *ApiResponse {
	return &ApiResponse{error, message, data}
}

package writer

import (
	"encoding/json"
	"net/http"
)

type ApiResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApiResponse(error bool, message string, data interface{}) *ApiResponse {
	return &ApiResponse{error, message, data}
}

type Writer interface {
	WriteResponse(w http.ResponseWriter, status int, response *ApiResponse)
}
type HttpWriter struct{}

func NewHttpWriter() *HttpWriter {
	return &HttpWriter{}
}

func (writer *HttpWriter) WriteResponse(w http.ResponseWriter, status int, response *ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

package writer

import (
	"encoding/json"
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/exceptions"
)

type Writer interface {
	WriteResponse(w http.ResponseWriter, status int, response *internal.ApiResponse)
	WriteError(w http.ResponseWriter, err error)
}
type HttpWriter struct {
	errorWriter *exceptions.ErrorWriter
}

func NewHttpWriter() *HttpWriter {
	return &HttpWriter{}
}

func (writer *HttpWriter) WriteResponse(w http.ResponseWriter, status int, response *internal.ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (writer *HttpWriter) WriteError(w http.ResponseWriter, err error) {
	writer.errorWriter.WriteError(w, err)
}

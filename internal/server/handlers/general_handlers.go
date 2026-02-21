package handlers

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
)

func PingHandler(hw *writer.HttpWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := throwError()
		if err != nil {
			hw.WriteError(w, err)
			return
		}
		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "pong", nil))
	}
}

func throwError() error {
	return exceptions.NewAppError("TEST_ERROR", "Test exception handling", nil)
}

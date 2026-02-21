package handlers

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal/server/writer"
)

func PingHandler(hw *writer.HttpWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hw.WriteResponse(w, http.StatusOK, writer.NewApiResponse(false, "pong", nil))
	}
}

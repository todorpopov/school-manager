package handlers

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/server/writer"
)

func PingHandler(hw *writer.HttpWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hw.WriteResponse(w, http.StatusOK, internal.NewApiResponse(false, "pong", nil))
		return
	}
}

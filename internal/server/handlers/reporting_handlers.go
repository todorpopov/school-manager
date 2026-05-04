package handlers

import (
	"net/http"

	"github.com/todorpopov/school-manager/internal"
	"github.com/todorpopov/school-manager/internal/reporting"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"go.uber.org/zap"
)

func ReportingQueryHandler(hw *writer.HttpWriter, rptSvc reporting.IDynamicReportingService, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, decodeErr := decodeRequestBodyInto[reporting.ReportQueryRequest](r, logger)
		if decodeErr != nil {
			hw.WriteError(w, decodeErr)
			return
		}

		data, err := rptSvc.GenerateReport(r.Context(), &request)
		if err != nil {
			logger.Error("Failed to run reporting query", zap.Error(err))
			hw.WriteError(w, err)
			return
		}

		logger.Info("Reporting query ran successfully", zap.Any("request", request))
		resp := internal.NewApiResponse(false, "Reporting query ran successfully", data)
		hw.WriteResponse(w, http.StatusCreated, resp)
	}
}

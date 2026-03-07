package middleware

import (
	"net/http"
	"time"

	"github.com/todorpopov/school-manager/internal/exceptions"
	"github.com/todorpopov/school-manager/internal/server/writer"
	"github.com/todorpopov/school-manager/internal/user_auth"
	"go.uber.org/zap"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func Logging(logger *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Info("Incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info("Request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", rw.statusCode),
				zap.Duration("duration", duration),
			)
		})
	}
}

func RequireRoles(hw *writer.HttpWriter, authSvc user_auth.IAuthService, roles ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionId := r.Header.Get("X-Session-Id")
			if sessionId == "" {
				hw.WriteError(w, exceptions.NewRequestValidationError("Session not provided"))
				return
			}

			request := user_auth.AuthRequest{SessionId: sessionId, RequiredRoles: roles}
			authorized, err := authSvc.IsRequestAuthorized(r.Context(), &request)
			if err != nil {
				hw.WriteError(w, err)
				return
			}

			if !authorized {
				hw.WriteError(w, exceptions.NewAppError("UNAUTHORIZED", "Unauthorized", nil))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

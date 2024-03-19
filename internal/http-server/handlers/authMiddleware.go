package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	invalidToken        = "Invalid Token"
	invalidAuthHeader   = "Invalid Authorization header"
	unauthorized        = "Unauthorized"
)

func (h *Handlers) authMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middleware_authMiddleware"

			header := r.Header.Get(authorizationHeader)
			if header == "" {
				logger.Error(op, slog.String("err", unauthorized))
				newErrorResponse(w, http.StatusUnauthorized, unauthorized)
				return
			}

			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 {
				logger.Error(op, slog.String("err", invalidAuthHeader))
				newErrorResponse(w, http.StatusBadRequest, invalidAuthHeader)
				return
			}
			pDataUser, err := h.service.Auth.ParseToken(headerParts[1])
			if err != nil {
				logger.Error(op, slog.String("err", invalidToken))
				newErrorResponse(w, http.StatusUnauthorized, invalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), parseDataUser, pDataUser)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

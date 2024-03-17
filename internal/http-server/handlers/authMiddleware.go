package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

//type parseDataUser struct {
//	id    int
//	email string
//	role  string
//}

func (h *Handlers) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
			return
		}
		fmt.Println("==========")
		pDataUser, err := h.service.Auth.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), parseDataUser, pDataUser)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

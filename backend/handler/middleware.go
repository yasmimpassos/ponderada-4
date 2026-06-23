package handler

import (
	"context"
	"net/http"
	"strings"

	"racha-historico/service"
)

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "token não fornecido")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := service.ValidateJWT(tokenString)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "token inválido")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

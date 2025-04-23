package middleware

import (
	"context"
	"net/http"
	"strings"

	http_handler "github.com/DarRo9/pvz_service/internal/handler"
	"github.com/DarRo9/pvz_service/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http_handler.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			http_handler.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

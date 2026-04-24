package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Geze296/orderhub/api-service/internal/auth"
	"github.com/Geze296/orderhub/api-service/internal/http/helper"
)

type contextKey string

const UserIdKey contextKey = "user_id"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				helper.WriteError(w, http.StatusUnauthorized, "Authorization Header required")
				return 
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2{
				helper.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return 
			}
			if parts[0] != "Bearer" {
				helper.WriteError(w, http.StatusUnauthorized, "Using wrong authorization method")
				return 
			}

			claims, err := auth.ParseToken(jwtSecret, parts[1])
			if err != nil {
				helper.WriteError(w, http.StatusUnauthorized, "Invalid Token")
				return 
			}

			fmt.Println("Claims user id: ", claims.UserId)
			ctx := context.WithValue(r.Context(), UserIdKey, int64(claims.UserId))
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIdFromContext(ctx context.Context) (int64, bool) {
	userId, ok := ctx.Value(UserIdKey).(int64)
	fmt.Println("user id: ",userId, "ok:", ok)
	return userId, ok
}
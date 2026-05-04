package middleware

import (
	"net/http"

	"github.com/Geze296/orderhub/api-service/internal/http/helper"
)

func RequireRole(requiredRole string) func(http.Handler) http.Handler{
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := UserRoleFromContext(r.Context())
			if !ok {
				helper.WriteError(w, http.StatusForbidden, "missing role")
				return 
			}

			if role != requiredRole {
				helper.WriteError(w, http.StatusForbidden, "wrong permission")
				return 
			}
			h.ServeHTTP(w, r)
		})
	}
}
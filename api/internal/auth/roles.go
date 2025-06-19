package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func RoleAuth(allowed ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			role, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "missing role", http.StatusUnauthorized)
				return
			}
			for _, a := range allowed {
				if a == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

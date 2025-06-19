package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth = jwtauth.New("HS256", secretKey, nil)

// Middleware verifica e autentica o JWT
func MustAuth(next http.Handler) http.Handler {
	return jwtauth.Authenticator(TokenAuth)(next)
}

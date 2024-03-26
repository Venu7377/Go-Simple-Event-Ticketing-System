package middleware

import (
	s "myProject/sessionManagement"
	"net/http"
)

func CombinedMiddleware(next http.Handler) http.Handler {
	return s.SessionMiddleware(BasicAuthMiddleware(next))
}

package middleware

import (
	"net/http"
)

var credentials = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type MiddlewareFunc func(http.Handler) http.Handler

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok {
			unauthorized(w)
			return
		}

		expectedPassword, userExists := credentials[username]

		if !userExists || password != expectedPassword {
			unauthorized(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("401 Unauthorized\n"))
}

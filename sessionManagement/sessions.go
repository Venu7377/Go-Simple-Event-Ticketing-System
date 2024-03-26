package sessionManagement

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

var credentials = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func init() {
	store = sessions.NewCookieStore([]byte("your-secret-key"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   180,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials_struct struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials_struct)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	expectedPassword, userExists := credentials[credentials_struct.Username]

	if !userExists || credentials_struct.Password != expectedPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = credentials_struct.Username
	session.Values["role"] = "user"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Login Successful!")
}

// type MiddlewareFunc func(http.Handler) http.Handler

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Unauthorized: Session expired or invalid", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		fmt.Println("error while ending session", err)
	}
	session.Values["authenticated"] = false
	session.Save(r, w)

	fmt.Fprintln(w, "Logout Successful!")
}

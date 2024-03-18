package main

import (
	mw "myProject/middleware"
	m "myProject/sessionManagement"
	"net/http"

	h "myProject/admin"

	db "myProject/db"

	"github.com/gorilla/mux"
)

func main() {
	db.InitializeMySQLDB()
	r := mux.NewRouter()
	r.HandleFunc("/login", m.LoginHandler)
	r.HandleFunc("/logout", m.LogoutHandler)
	r.Use(m.SessionMiddleware)
	r.Use(mw.BasicAuthMiddleware)
	r.HandleFunc("/add", h.AddEventHandler)
	r.HandleFunc("/getAll", h.GetAllEventsHandler)
	r.HandleFunc("/update", h.UpdateEventHandler)
	r.HandleFunc("/delete", h.DeleteEventHandler)

	wrappedmux := mw.NewLogger(r)

	http.ListenAndServe(":80", wrappedmux)

}

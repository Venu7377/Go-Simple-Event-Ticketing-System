package main

import (
	"fmt"
	"log"
	h "myProject/admin"
	mw "myProject/middleware"
	m "myProject/sessionManagement"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	go func() {
		if err := h.SyncCachePeriodically(); err != nil {
			log.Printf("Error in cache sync: %v", err)
		}
	}()
	r := mux.NewRouter()
	r.HandleFunc("/login", m.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", m.LogoutHandler).Methods("GET")
	r.HandleFunc("/bookTicket", h.BookTicketsHandler).Methods("POST")

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(mw.CombinedMiddleware)

	adminRouter.HandleFunc("/add", h.AddEventHandler).Methods("POST")
	adminRouter.HandleFunc("/getAllUsingRedis", h.GetAllEventsHandler).Methods("GET")
	adminRouter.HandleFunc("/get/{id}", h.GetEventByIdHandler).Methods("GET")
	adminRouter.HandleFunc("/update/{id}", h.UpdateEventHandler).Methods("PUT")
	adminRouter.HandleFunc("/delete/{id}", h.DeleteEventHandler).Methods("DELETE")
	wrappedmux := mw.NewLogger(r)

	fmt.Println("Server is Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", wrappedmux))

}

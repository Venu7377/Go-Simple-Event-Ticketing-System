package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
)

var repo = &Repository{}

func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	repo.GetAllEvents(w)
}

func AddEventHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}
	repo.AddEvent(w, &event)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}
	repo.UpdateEvent(w, &event)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	repo.DeleteEvent(w, id)
}

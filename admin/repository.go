package admin

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetAllEvents(w http.ResponseWriter) {
	var events []*Event
	if err := r.DB.Find(&events).Error; err != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "Failed to encode events", http.StatusInternalServerError)
	}
}

func (r *Repository) AddEvent(w http.ResponseWriter, event *Event) {
	if err := r.DB.Create(event).Error; err != nil {
		http.Error(w, "Failed to add event", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}
}

func (r *Repository) UpdateEvent(w http.ResponseWriter, event *Event) {
	if err := r.DB.Save(event).Error; err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}

}

func (r *Repository) DeleteEvent(w http.ResponseWriter, id uint64) {
	if err := r.DB.Delete(&Event{}, id).Error; err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Event deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

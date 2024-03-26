package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	a "myProject/model"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Repository struct {
	DB     *gorm.DB
	Redis  *redis.Client
	ctx    context.Context
	expiry time.Duration
}

func NewRepository(db *gorm.DB, redis *redis.Client) *Repository {
	return &Repository{
		DB:     db,
		Redis:  redis,
		ctx:    context.Background(),
		expiry: time.Minute * 5,
	}
}

func (r *Repository) GetAllEvents(w http.ResponseWriter) {
	eventsJSON, err := r.Redis.Get(r.ctx, "events").Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(eventsJSON))
		return
	}

}

func (r *Repository) AddEvent(w http.ResponseWriter, event *a.Event) {
	if err := r.DB.Create(event).Error; err != nil {
		http.Error(w, "Failed to add event", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}
}

func (r *Repository) UpdateEvent(w http.ResponseWriter, event *a.Event, id int) {
	existingEvent := &a.Event{}
	if err := r.DB.First(existingEvent, id).Error; err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}
	if event.Name != "" {
		existingEvent.Name = event.Name
	}
	if event.Date != "" {
		existingEvent.Date = event.Date
	}
	if event.Place != "" {
		existingEvent.Place = event.Place
	}
	if event.NumberOfTickets != 0 {
		mutex.Lock()
		defer mutex.Unlock()
		existingEvent.NumberOfTickets = event.NumberOfTickets
	}

	if err := r.DB.Save(existingEvent).Error; err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(existingEvent); err != nil {
		http.Error(w, "Failed to encode event data", http.StatusInternalServerError)
		return
	}
}

func (r *Repository) DeleteEvent(w http.ResponseWriter, id int) {
	if err := r.DB.Delete(&a.Event{}, id).Error; err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Event deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func (r *Repository) FindEventByID(w http.ResponseWriter, id int) {
	var event a.Event
	if err := r.DB.Where("id = ?", id).First(&event).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, fmt.Sprintf("Event with ID %d not found", id), http.StatusNotFound)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "Failed to encode event data", http.StatusInternalServerError)
		return
	}
}

func (r *Repository) syncCacheWithDatabase() error {
	logFilePath := "./logs/app.log"
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	var events []*a.Event
	if err := r.DB.Find(&events).Error; err != nil {
		log.Printf("Failed to fetch events:%v", err)
		return err
	}

	events_JSON, err := json.Marshal(events)
	if err != nil {
		log.Printf("Failed to encode events:%v", err)
		return err
	}
	if err := r.Redis.Set(r.ctx, "events", events_JSON, r.expiry).Err(); err != nil {
		log.Printf("Failed to cache events:%v", err)
		return err
	}
	return nil
}

func SyncCachePeriodically() error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	var err error
	for {
		<-ticker.C
		err = repo.syncCacheWithDatabase()
		if err != nil {
			return err
		}

	}
}

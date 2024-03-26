package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"myProject/db"
	a "myProject/model"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var repo *Repository
var dbcon *gorm.DB
var mutex = &sync.Mutex{}

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	dbcon, _ = db.InitializeMySQLDB()
	repo = NewRepository(dbcon, redisClient)
}

func GetAllEventsHandler(w http.ResponseWriter, r *http.Request) {
	repo.GetAllEvents(w)
}
func AddEventHandler(w http.ResponseWriter, r *http.Request) {
	var event a.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Failed to encode event", http.StatusInternalServerError)
	}
	repo.AddEvent(w, &event)
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var event a.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Failed to decode event", http.StatusInternalServerError)
	}
	repo.UpdateEvent(w, &event, id)
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	repo.DeleteEvent(w, id)
}

func GetEventByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	repo.FindEventByID(w, id)
}

// ------------------------------------------------------------

func BookTicketsHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		EventID      int    `json:"event_id"`
		UserEmail    string `json:"user_email"`
		NumOfTickets int32  `json:"num_of_tickets"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	eventID := requestBody.EventID
	userEmail := requestBody.UserEmail
	numOfTickets := requestBody.NumOfTickets

	if eventID == 0 || userEmail == "" || numOfTickets == 0 {
		http.Error(w, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	if !strings.Contains(userEmail, "@") {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	existingEvent := &a.Event{}
	if err := dbcon.First(existingEvent, eventID).Error; err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	if numOfTickets < 0 {
		http.Error(w, "Enter Valid Number of Tickets", http.StatusBadRequest)
		return
	}

	if existingEvent.NumberOfTickets < numOfTickets {
		http.Error(w, "Not enough tickets available", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	success := bookTickets(*existingEvent, userEmail, numOfTickets)
	if success {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully booked %d tickets for event %d", numOfTickets, eventID)
	} else {
		http.Error(w, "Failed to book tickets", http.StatusInternalServerError)
	}
}

func bookTickets(event a.Event, userEmail string, numOfTickets int32) bool {
	event.NumberOfTickets -= numOfTickets
	result := dbcon.Save(&event)
	if result.Error != nil {
		return false
	}
	go sendTicket(numOfTickets, userEmail)

	return true
}

func sendTicket(userTickets int32, email string) {
	time.Sleep(10 * time.Second)
	fmt.Println("######################")
	fmt.Printf("Sending ticket:\n %v tickets \nto email address %v\n", userTickets, email)
	fmt.Println("######################")
}

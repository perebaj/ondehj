package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/perebaj/ondehj/event"
)

const (
	eventPath   = "/event"
	eventPathId = "/event/{id}"
)

func deleteEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("deleteEventHandler")
		if r.Method != http.MethodDelete {
			fmt.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid id", http.StatusBadRequest)
		}

		fmt.Println("Deleting event with id: ", idStr)
		err = eventRepo.Delete(r.Context(), id)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}
	}
	return http.HandlerFunc(fn)
}

func postCreateEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("postCreateEventHandler")

		if r.Method != http.MethodPost {
			fmt.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var event event.Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("Creating event")
		_, err = eventRepo.Create(r.Context(), event)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(event)

	}
	return http.HandlerFunc(fn)
}

func getAllEventsHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getAllEventsHandler")
		fmt.Println("Get all events")
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		events, err := eventRepo.All(r.Context())
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(events)
	}
	return http.HandlerFunc(fn)
}

func getByIDHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("getByIDHandler")
		fmt.Println("Get event by id")
		if r.Method != http.MethodGet {
			fmt.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid id", http.StatusBadRequest)
		}

		event, err := eventRepo.GetByID(r.Context(), id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		fmt.Println(event)
	}
	return http.HandlerFunc(fn)
}

func HandlerFactory(db *pgxpool.Pool) http.Handler {
	//Group all handler of the API and return a http.Handler
	router := mux.NewRouter()
	eventSQLRepo := event.EventSQLRepository(db)

	//event
	router.HandleFunc(eventPath, getAllEventsHandler(eventSQLRepo)).Methods(http.MethodGet)
	router.HandleFunc(eventPath, postCreateEventHandler(eventSQLRepo)).Methods(http.MethodPost)
	router.HandleFunc(eventPathId, deleteEventHandler(eventSQLRepo)).Methods(http.MethodDelete)
	router.HandleFunc(eventPathId, getByIDHandler(eventSQLRepo)).Methods(http.MethodGet)

	return router
}

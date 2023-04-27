package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/perebaj/ondehj/event"
	"golang.org/x/exp/slog"
)

const (
	eventPath   = "/event"
	eventPathId = "/event/{id}"
)

func deleteEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Calling deleteEventHandler")
		if r.Method != http.MethodDelete {
			slog.Error("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error(fmt.Sprintf("Invalid id: %s", idStr), "error", err)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		slog.Info(fmt.Sprintf("Getting event with id: %d", id))
		_, err = eventRepo.GetByID(r.Context(), id)
		if err != nil {
			slog.Error("Event doesn't exist", "error", err)
			http.Error(w, "Event doesn't exist", http.StatusNotFound)
			return
		}

		slog.Info(fmt.Sprintf("Deleting event with id: %d", id))
		err = eventRepo.Delete(r.Context(), id)

		if err != nil {
			slog.Error("Delete failed", "error", err)
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}
		slog.Info("Event deleted successfully")
	}
	return http.HandlerFunc(fn)
}

func postCreateEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("postCreateEventHandler")

		if r.Method != http.MethodPost {
			slog.Error("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var requestEvent event.Event
		err := json.NewDecoder(r.Body).Decode(&requestEvent)
		if err != nil {
			slog.Error("Error decoding event", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if requestEvent == (event.Event{}) {
			slog.Error("Empty event")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		slog.Info("Creating event")
		createdEvent, err := eventRepo.Create(r.Context(), requestEvent)
		if err != nil {
			slog.Error("Error creating new Event", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		eventJson, err := json.Marshal(createdEvent)
		if err != nil {
			slog.Error("Error marshalling events", "error", err)
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(eventJson)
		slog.Info("Event created successfully")

	}
	return http.HandlerFunc(fn)
}

func getAllEventsHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getAllEventsHandler")
		if r.Method != http.MethodGet {
			slog.Error("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		events, err := eventRepo.All(r.Context())
		if err != nil {
			slog.Error("Error getting all events", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		eventJson, err := json.Marshal(events)
		if err != nil {
			slog.Error("Error marshalling events", "error", err)
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(eventJson)
		slog.Info("Events retrieved successfully")
	}
	return http.HandlerFunc(fn)
}

func getByIDHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getByIDHandler")
		if r.Method != http.MethodGet {
			slog.Error("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error(fmt.Sprintf("Invalid id: %s", idStr), "error", err)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		event, err := eventRepo.GetByID(r.Context(), id)
		if err != nil {
			slog.Error("Event not found", "error", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		eventJson, err := json.Marshal(event)
		if err != nil {
			slog.Error("Error marshalling events", "error", err)
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(eventJson)
		slog.Info("Event retrieved successfully")
	}
	return http.HandlerFunc(fn)
}

func Update(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Update")
		if r.Method != http.MethodPut {
			slog.Error("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var newEvent event.Event
		err := json.NewDecoder(r.Body).Decode(&newEvent)
		if err != nil {
			slog.Error("Error decoding event", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error(fmt.Sprintf("Invalid id: %s", idStr), "error", err)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		_, err = eventRepo.GetByID(r.Context(), id)
		if err != nil {
			slog.Error("Event not found", "error", err)
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		updatedEvent, err := eventRepo.Update(r.Context(), id, newEvent)
		if err != nil {
			slog.Error("Update failed", "error", err)
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}
		updatedEventJson, err := json.Marshal(updatedEvent)
		if err != nil {
			slog.Error("Error marshalling events", "error", err)
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(updatedEventJson)
		slog.Info("Event updated successfully")
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
	router.HandleFunc(eventPathId, Update(eventSQLRepo)).Methods(http.MethodPut)

	return router
}

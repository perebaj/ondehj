package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/httplog"
	"github.com/go-openapi/runtime/middleware"
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
		log := httplog.LogEntry(r.Context())
		log.Info().Msg("Calling deleteEventHandler")
		if r.Method != http.MethodDelete {
			log.Error().Msg("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Err(err).Msgf("Invalid id: %s", idStr)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		log.Info().Msgf("Getting event with id: %d", id)
		_, err = eventRepo.GetByID(r.Context(), id, log)
		if err != nil {
			log.Err(err).Msg("Event doesn't exist")
			http.Error(w, "Event doesn't exist", http.StatusNotFound)
			return
		}

		log.Info().Msgf("Deleting event with id: %d", id)
		err = eventRepo.Delete(r.Context(), id, log)

		if err != nil {
			log.Err(err).Msg("Delete failed")
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}
		log.Info().Msg("Event deleted successfully")
	}
	return http.HandlerFunc(fn)
}

func postCreateEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := httplog.LogEntry(r.Context())
		log.Info().Msg("postCreateEventHandler")

		if r.Method != http.MethodPost {
			log.Error().Msg("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Decode the request body into a Event struct
		var requestEvent event.Event
		err := json.NewDecoder(r.Body).Decode(&requestEvent)
		if err != nil {
			log.Err(err).Msg("Error decoding event")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Validate the request body
		if requestEvent == (event.Event{}) || requestEvent.Title == "" {
			// If event is empty or title is empty, return an error
			log.Error().Msg("Invalid Event")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info().Msg("Creating event")
		createdEvent, err := eventRepo.Create(r.Context(), requestEvent, log)
		if err != nil {
			log.Err(err).Msg("Error creating new Event")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		eventJson, err := json.Marshal(createdEvent)
		if err != nil {
			log.Err(err).Msg("Error marshalling events")
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(eventJson)
		log.Info().Msg("Event created successfully")

	}
	return http.HandlerFunc(fn)
}

func getAllEventsHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := httplog.LogEntry(r.Context())
		log.Info().Msg("getAllEventsHandler")
		if r.Method != http.MethodGet {
			log.Error().Msg("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		events, err := eventRepo.All(r.Context(), log)
		if err != nil {
			log.Err(err).Msg("Error retrieving events")
			log.Error().AnErr("error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		eventJson, err := json.Marshal(events)
		if err != nil {
			log.Err(err).Msg("Error marshalling events")
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(eventJson)
		log.Info().Msg("Events retrieved successfully")
	}
	return http.HandlerFunc(fn)
}

func getByIDHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := httplog.LogEntry(r.Context())
		log.Info().Msg("getByIDHandler")
		if r.Method != http.MethodGet {
			log.Error().Msg("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Err(err).Msgf("Invalid id: %s", idStr)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		event, err := eventRepo.GetByID(r.Context(), id, log)
		if err != nil {
			log.Err(err).Msg("Event not found")
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		eventJson, err := json.Marshal(event)
		if err != nil {
			log.Err(err).Msg("Error marshalling events")
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(eventJson)
		log.Info().Msg("Event retrieved successfully")
	}
	return http.HandlerFunc(fn)
}

func Update(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log := httplog.LogEntry(r.Context())
		log.Info().Msg("Update")
		if r.Method != http.MethodPut {
			log.Error().Msg("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var newEvent event.Event
		err := json.NewDecoder(r.Body).Decode(&newEvent)
		if err != nil {
			log.Err(err).Msg("Error decoding event")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		idStr := mux.Vars(r)["id"]
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Err(err).Msgf("Invalid id: %s", idStr)
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		_, err = eventRepo.GetByID(r.Context(), id, log)
		if err != nil {
			log.Err(err).Msg("Event not found")
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}
		updatedEvent, err := eventRepo.Update(r.Context(), id, newEvent, log)
		if err != nil {
			log.Err(err).Msg("Update failed")
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}
		updatedEventJson, err := json.Marshal(updatedEvent)
		if err != nil {
			log.Err(err).Msg("Error marshalling events")
			http.Error(w, "Error marshalling events", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(updatedEventJson)
		log.Info().Msg("Event updated successfully")
	}
	return http.HandlerFunc(fn)
}

func HandlerFactory(db *pgxpool.Pool) http.Handler {
	//Group all handler of the API and return a http.Handler
	//structured logs
	logger := httplog.NewLogger("http", httplog.Options{
		JSON:     false,
		LogLevel: "info",
		Concise:  true,
	})
	httpLogMiddleware := httplog.RequestLogger(logger)
	router := mux.NewRouter()
	eventSQLRepo := event.EventSQLRepository(db)

	//event
	router.Use(httpLogMiddleware)
	router.HandleFunc(eventPath, getAllEventsHandler(eventSQLRepo)).Methods(http.MethodGet)
	router.HandleFunc(eventPath, postCreateEventHandler(eventSQLRepo)).Methods(http.MethodPost)
	router.HandleFunc(eventPathId, deleteEventHandler(eventSQLRepo)).Methods(http.MethodDelete)
	router.HandleFunc(eventPathId, getByIDHandler(eventSQLRepo)).Methods(http.MethodGet)
	router.HandleFunc(eventPathId, Update(eventSQLRepo)).Methods(http.MethodPut)
	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "openapi.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)
	router.Handle("/openapi.yaml", http.FileServer(http.Dir("./")))

	return router
}

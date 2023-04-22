package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/perebaj/ondehj/event"
)

const (
	eventPath = "/event"
)

func eventHandler(eventRepo event.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Println("Event GET handler")
			getAllEventsHandler(eventRepo)(w, r)
		case http.MethodPost:
			fmt.Println("Event POST handler")
			postCreateEventHandler(eventRepo)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

}

func postCreateEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("postCreateEventHandler")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
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
			w.WriteHeader(http.StatusMethodNotAllowed)
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

func HandlerFactory(db *pgxpool.Pool) http.Handler {
	//Group all handler of the API and return a http.Handler
	mux := http.NewServeMux()
	eventSQLRepo := event.EventSQLRepository(db)

	eventHandler := eventHandler(eventSQLRepo)
	mux.Handle(eventPath, eventHandler)
	return mux
}

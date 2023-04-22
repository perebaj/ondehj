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

func createEventHandler(eventRepo event.Repository) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start handler function with parameter")

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

func HandlerFactory(db *pgxpool.Pool) http.Handler {
	//Group all handler of the API and return a http.Handler
	mux := http.NewServeMux()
	eventSQLRepo := event.EventSQLRepository(db)

	createEventHandler := createEventHandler(eventSQLRepo)
	mux.Handle(eventPath, createEventHandler)
	return mux
}

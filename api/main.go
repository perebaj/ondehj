package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/perebaj/ondehj/event"
)

func createEvent(w http.ResponseWriter, r *http.Request) {
	var event event.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(event)

}

func main() {
	mux := http.NewServeMux()

	handlerCreateEvent := http.HandlerFunc(createEvent)

	mux.Handle("/event", handlerCreateEvent)

	log.Print("Listening...")
	http.ListenAndServe(":3000", mux)
}

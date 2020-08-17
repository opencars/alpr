package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

// The Handler helps to handle errors in one place.
type Handler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h(w, r); err != nil {
		switch e := err.(type) {
		case Error:
			// We can retrieve the status here and write out a specific HTTP status code.
			log.Printf("HTTP %d - %s\n", e.Status(), e)
			w.WriteHeader(e.Status())
			if err := json.NewEncoder(w).Encode(e); err != nil {
				panic(err)
			}
		default:
			// Any error types we don't specifically look out for default to serving a HTTP 500
			log.Printf("HTTP %d - %s\n", http.StatusInternalServerError, e)
			w.WriteHeader(http.StatusInternalServerError)
			e = NewError(http.StatusInternalServerError, "Something went wrong")
			if err := json.NewEncoder(w).Encode(e); err != nil {
				panic(err)
			}
		}
	}
}

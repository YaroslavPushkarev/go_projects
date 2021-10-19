package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (j JokesHandler) GetId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var jokes Joke

	id := r.URL.Query().Get("id")

	query := map[string]interface{}{"id": id}

	err := FindId(j.collection, query).Decode(&jokes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Midl(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("hello")
		next(w, r)
	}
}

package api

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (j JokesHandler) Funniest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination, err := j.ParseSkipAndLimit(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cursor, err := j.Storage.FunniestJokes(bson.D{}, pagination.Limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(cursor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

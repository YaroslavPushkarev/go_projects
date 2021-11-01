package api

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (j JokesHandler) RandomJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination, err := j.ParseSkipAndLimit(w, r)

	cursor, err := j.Storage.Random(bson.M{}, pagination.Limit)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err = json.NewEncoder(w).Encode(cursor); err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
	}
}

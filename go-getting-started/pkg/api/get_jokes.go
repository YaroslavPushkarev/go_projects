package api

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (j JokesHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination, err := j.ParseSkipAndLimit(w, r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jokes, err := j.Storage.GetJokes(bson.M{}, pagination.Limit, pagination.Skip)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

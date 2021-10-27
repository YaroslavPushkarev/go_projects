package api

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (j JokesHandler) GetJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jokes, err := j.Storage.GetJokes(bson.M{})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
	}
}

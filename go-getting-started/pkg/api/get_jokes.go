package api

import (
	"encoding/json"
	"net/http"

	storage "github.com/heroku/go-getting-started/pkg/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func GetJokes(db storage.JokesInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		jokes, err := db.GetJokes(bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(jokes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
		}
	}
}

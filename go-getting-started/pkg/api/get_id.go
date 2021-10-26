package api

import (
	"encoding/json"
	"net/http"

	storage "github.com/heroku/go-getting-started/pkg/storage/mongo"
)

func GetId(db storage.JokesInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("content-type", "application/json")

		id := r.URL.Query().Get("id")

		res, err := db.FindId(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNoContent)
		}
	}

}

package api

import (
	"encoding/json"
	"net/http"
)

func (j JokesHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	search := r.URL.Query().Get("search")

	cursor, err := j.Storage.FindJoke(search)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(cursor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNoContent)
	}
}

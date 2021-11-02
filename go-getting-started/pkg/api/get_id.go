package api

import (
	"encoding/json"
	"net/http"
)

func (j JokesHandler) GetId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id := r.URL.Query().Get("id")

	res, err := j.Storage.FindId(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

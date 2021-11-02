package api

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strconv"

	"github.com/heroku/go-getting-started/pkg/models"
)

func (j JokesHandler) CreateJoke(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body := r.URL.Query().Get("body")

	score, err := strconv.Atoi(r.URL.Query().Get("score"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	title := r.URL.Query().Get("title")

	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	randomID := fmt.Sprintf("%X", b)

	err = j.Storage.InsertJoke(models.Joke{Body: body, ID: randomID, Score: score, Title: title})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

}

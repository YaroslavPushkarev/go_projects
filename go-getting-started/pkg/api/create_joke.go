package api

import (
	"fmt"
	"net/http"

	"github.com/heroku/go-getting-started/pkg/models"
)

func (j JokesHandler) CreateJoke(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := j.Storage.InsertJoke(models.Joke{Body: "They're both Paris sites", ID: "2ds4s", Score: 3, Title: "What do a tick and the Eiffel Tower have in common?"})
	if err != nil {
		fmt.Println(err)
	}

}

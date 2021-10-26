package api

import (
	"fmt"
	"net/http"

	"github.com/heroku/go-getting-started/pkg/models"
	storage "github.com/heroku/go-getting-started/pkg/storage/mongo"
)

func CreateJoke(db storage.JokesInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		joke := models.Joke{Body: "They're both Paris sites", ID: "2ds4s", Score: 3, Title: "What do a tick and the Eiffel Tower have in common?"}

		res, err := db.InsertJoke(joke)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res)
	}
}

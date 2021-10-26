package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/pkg/api"
	"github.com/heroku/go-getting-started/pkg/models"
	storage "github.com/heroku/go-getting-started/pkg/storage/mongo"
)

const port = ":8080"

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []models.Joke{}
	err := json.Unmarshal(content, &jokes)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.TODO()

	col := storage.ConnectDB(os.Getenv("MONGODB_URI"))

	client := &storage.JokesHandler{
		Collection: col,
		Ctx:        ctx,
	}

	http.HandleFunc("/jokes/id", api.GetId(client))
	http.HandleFunc("/jokes/jokes", api.GetJokes(client))
	http.HandleFunc("/jokes/create", api.CreateJoke(client))

	log.Fatal(http.ListenAndServe(port, nil))
}

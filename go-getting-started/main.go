package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/heroku/go-getting-started/pkg/api"
	storage "github.com/heroku/go-getting-started/pkg/storage/mongo"
)

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []api.Joke{}
	err := json.Unmarshal(content, &jokes)
	if err != nil {
		fmt.Println(err)
	}

	collection := storage.ConnectDB(os.Getenv("MONGODB_URI"))
	jh := api.JokesHandler{jokes, collection}

	http.Handle("/j", jh)
	http.HandleFunc("/jokes", jh.getJokes)
	http.HandleFunc("/jokes/id", jh.getId)
	http.HandleFunc("/jokes/search", jh.search)
	http.HandleFunc("/jokes/funniest", jh.funniest)
	http.HandleFunc("/jokes/random", jh.randomJokes)
	http.HandleFunc("/jokes/create", jh.createJoke)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

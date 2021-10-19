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

const port = ":8080"

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []api.Joke{}
	err := json.Unmarshal(content, &jokes)
	if err != nil {
		fmt.Println(err)
	}

	_ = storage.ConnectDB(os.Getenv("MONGODB_URI"))

	jh := api.JokesHandler{}

	http.HandleFunc("/jokes/id", api.Midl(jh.GetId))
	log.Fatal(http.ListenAndServe(port, nil))
}

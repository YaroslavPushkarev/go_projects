package main

import (
	"log"
	"net/http"

	"github.com/heroku/go-getting-started/pkg/api"

	"github.com/heroku/go-getting-started/pkg/storage/mongo"
)

const port = ":8080"

func main() {

	// collection := storage.ConnectDB(os.Getenv("MONGODB_URI"))
	collection := mongo.ConnectDB("mongodb+srv://jokesdb:jokesdb@joke.kxki9.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	storage := mongo.JokesStorage{
		Collection: collection,
	}

	client := api.JokesHandler{
		Storage: storage,
	}

	http.HandleFunc("/jokes/create", client.CreateJoke)
	http.HandleFunc("/jokes/jokes", client.GetJokes)
	http.HandleFunc("/jokes/id", client.GetId)
	http.HandleFunc("/jokes/search", client.Search)

	log.Fatal(http.ListenAndServe(port, nil))
}

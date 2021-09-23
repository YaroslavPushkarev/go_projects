package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/heroku/go-getting-started/controllers"

	"github.com/heroku/go-getting-started/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Skip  int
	Limit int
}

type Joke struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Score int    `json:"score" bson:"score"`
	Body  string `json:"body" bson:"body"`
}

type jokesHandler struct {
	jokes []Joke
}

func (j jokesHandler) parseSkipAndLimit(w http.ResponseWriter, r *http.Request) (Pagination, error) {
	leng := len(j.jokes)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		limit = 1
	}
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		skip = 1
	}
	if skip > leng {
		w.WriteHeader(http.StatusBadRequest)
		return Pagination{}, nil
	}
	if skip < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return Pagination{}, nil
	}

	if limit > leng {
		limit = leng - skip
	}

	if limit < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return Pagination{}, nil
	}

	pagination := Pagination{Skip: skip, Limit: limit}
	return pagination, nil
}

func (j jokesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(j.jokes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		err := json.NewEncoder(w).Encode(j.jokes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if len(j.jokes) == 1 {
		err := json.NewEncoder(w).Encode(j.jokes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	jokes := j.jokes

	pagination, err := j.parseSkipAndLimit(w, r)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var collection = config.ConnectDB()

func (j jokesHandler) getJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var jokes []Joke

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pagination.Limit)).SetSkip(int64(pagination.Skip))

	cursor, err := collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var joke Joke
		err := cursor.Decode(&joke)
		if err != nil {
			panic(err)
		}
		jokes = append(jokes, joke)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) getId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var jokes Joke

	id := r.URL.Query().Get("id")

	query := map[string]interface{}{"id": id}

	err := collection.FindOne(context.TODO(), query).Decode(&jokes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) randomJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jokes []Joke

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		fmt.Println(err)
	}

	pipeline := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: pagination.Limit}}}}

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{pipeline})
	if err != nil {
		fmt.Println(err)
	}

	if err = cursor.All(context.TODO(), &jokes); err != nil {
		panic(err)
	}

	if err = json.NewEncoder(w).Encode(jokes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) funniest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jokes []Joke

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		fmt.Println(err)
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: -1}}).SetLimit(int64(pagination.Limit))

	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = cursor.All(context.TODO(), &jokes); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	search := r.URL.Query().Get("search")

	var jokes []Joke

	query := bson.M{
		"title": bson.M{
			"$regex": primitive.Regex{
				Pattern: search,
				Options: "i",
			},
		},
	}

	cursor, err := collection.Find(context.TODO(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for cursor.Next(context.TODO()) {
		var joke Joke
		err := cursor.Decode(&joke)
		if err != nil {
			panic(err)
		}
		jokes = append(jokes, joke)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []Joke{}
	err := json.Unmarshal(content, &jokes)
	if err != nil {
		fmt.Println(err)
	}

	jh := jokesHandler{jokes}

	http.Handle("/j", jh)
	http.HandleFunc("/jokes", jh.getJokes)
	http.HandleFunc("/jokes/id", jh.getId)
	http.HandleFunc("/jokes/search", search)
	http.HandleFunc("/jokes/funniest", jh.funniest)
	http.HandleFunc("/jokes/random", jh.randomJokes)
	http.HandleFunc("/jokes/create", controllers.CreateJoke)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

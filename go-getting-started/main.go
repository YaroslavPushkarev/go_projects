package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

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

var client *mongo.Client

func getId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id := r.URL.Query().Get("id")

	var jokes Joke

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := map[string]interface{}{"id": id}

	err := collection.FindOne(ctx, query).Decode(&jokes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func randomJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var jokes []Joke

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pip := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 5}}}}
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{pip})
	if err != nil {
		fmt.Println(err)
	}
	if err = cursor.All(ctx, &jokes); err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func funniest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jokes []Joke

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: 1}}).SetLimit(20)

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = cursor.All(ctx, &jokes); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jokes []Joke

	search := r.URL.Query().Get("search")

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := bson.M{
		"title": bson.M{
			"$regex": primitive.Regex{
				Pattern: search,
				Options: "i",
			},
		},
	}

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for cursor.Next(ctx) {
		var joke Joke
		err := cursor.Decode(&joke)
		if err != nil {
			panic(err)
		}
		jokes = append(jokes, joke)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []Joke{}
	err := json.Unmarshal(content, &jokes)
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://jokesdb:jokesdb@joke.kxki9.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)

	http.Handle("/jokes", jokesHandler{jokes})
	http.HandleFunc("/jokesdb", getId)
	http.HandleFunc("/jokes/search", search)
	http.HandleFunc("/jokes/funniest", funniest)
	http.HandleFunc("/jokes/random", randomJokes)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

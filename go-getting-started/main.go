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

func (j jokesHandler) getJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetLimit(int64(pagination.Limit)).SetSkip(int64(pagination.Skip))

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var joke Joke
		err := cursor.Decode(&joke)
		if err != nil {
			panic(err)
		}
		j.jokes = append(j.jokes, joke)
	}

	err = json.NewEncoder(w).Encode(j.jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) getId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	id := r.URL.Query().Get("id")

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	query := map[string]interface{}{"id": id}

	err := collection.FindOne(ctx, query).Decode(&j.jokes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(j.jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) randomJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		fmt.Println(err)
	}

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pipeline := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: pagination.Limit}}}}

	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{pipeline})
	if err != nil {
		fmt.Println(err)
	}

	if err = cursor.All(ctx, &j.jokes); err != nil {
		panic(err)
	}

	if err = json.NewEncoder(w).Encode(j.jokes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) funniest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database("Jokes").Collection("jokes")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: -1}}).SetLimit(20)

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = cursor.All(ctx, &j.jokes); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(j.jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (j jokesHandler) search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
		j.jokes = append(j.jokes, joke)
	}

	err = json.NewEncoder(w).Encode(j.jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func (j jokesHandler)createJoke(w http.ResponseWriter, r *http.Request) {
// 	collection := client.Database("Jokes").Collection("jokes")
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	_, err := collection.InsertOne(ctx, bson.D{
// 		{Key: "body", Value: "An im-pasta"},
// 		{Key: "id", Value: "aw42r54t"},
// 		{Key: "score", Value: 3},
// 		{Key: "title", Value: "What do you call a fake noodle?"},
// 	})

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []Joke{}
	err := json.Unmarshal(content, &jokes)
	if err != nil {
		fmt.Println(err)
	}

	jh := jokesHandler{jokes}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb+srv://jokesdb:jokesdb@joke.kxki9.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
	client, _ = mongo.Connect(ctx, clientOptions)

	http.Handle("/j", jh)
	http.HandleFunc("/jokes", jh.getJokes)
	http.HandleFunc("/jokes/id", jh.getId)
	http.HandleFunc("/jokes/search", jh.search)
	http.HandleFunc("/jokes/funniest", jh.funniest)
	http.HandleFunc("/jokes/random", jokesHandler{jokes}.randomJokes)
	// http.HandleFunc("/jokes/create", createJoke)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (j JokesHandler) getJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var jokes []Joke

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(pagination.Limit)).SetSkip(int64(pagination.Skip))

	cursor, err := j.collection.Find(context.TODO(), bson.M{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
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
		http.Error(w, err.Error(), http.StatusNoContent)
	}
}

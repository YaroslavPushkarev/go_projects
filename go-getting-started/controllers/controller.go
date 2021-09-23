package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/heroku/go-getting-started/config"
	"go.mongodb.org/mongo-driver/bson"
)

type Joke struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Score int    `json:"score" bson:"score"`
	Body  string `json:"body" bson:"body"`
}

var collection = config.ConnectDB()

func CreateJoke(w http.ResponseWriter, r *http.Request) {
	_, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "body", Value: "A im-pasta"},
		{Key: "id", Value: "aw42r54t"},
		{Key: "score", Value: 3},
		{Key: "title", Value: "What do you call a fake noodle?"},
	})

	if err != nil {
		fmt.Println(err)
	}
}

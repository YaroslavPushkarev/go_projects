package controllers

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Joke struct {
	ID    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Score int    `json:"score" bson:"score"`
	Body  string `json:"body" bson:"body"`
}

func InsertData(collection *mongo.Collection, joke Joke) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(context.Background(), joke)

	if err != nil {
		return res, err
	}
	return res, err
}

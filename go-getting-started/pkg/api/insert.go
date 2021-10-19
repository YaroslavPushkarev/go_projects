package api

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func InsertJoke(collection *mongo.Collection, joke Joke) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(context.Background(), joke)

	if err != nil {
		return res, err
	}

	return res, err
}

package controllers

import (
	"context"
	"fmt"

	"github.com/heroku/go-getting-started/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertData(collection *mongo.Collection, joke models.Joke) (*mongo.InsertOneResult, error) {
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		fmt.Println(err)
	}

	res, err := collection.InsertOne(context.Background(), joke)

	if err != nil {
		return res, err
	}

	return res, err
}

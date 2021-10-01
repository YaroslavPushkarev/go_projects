package controllers

import (
	"context"

	"github.com/heroku/go-getting-started/dbinterface"
	"github.com/heroku/go-getting-started/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertData(collection dbinterface.CollectionAPI, joke models.Joke) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(context.Background(), joke)

	if err != nil {
		return res, err
	}

	return res, err
}

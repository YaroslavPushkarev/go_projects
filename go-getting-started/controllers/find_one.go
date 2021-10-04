package controllers

import (
	"context"

	"github.com/heroku/go-getting-started/dbinterface"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindId(collection dbinterface.CollectionAPI, query map[string]interface{}) *mongo.SingleResult {
	res := collection.FindOne(context.TODO(), query)
	return res
}

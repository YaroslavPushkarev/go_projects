package api

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

func FindId(collection *mongo.Collection, query map[string]interface{}) *mongo.SingleResult {
	res := collection.FindOne(context.TODO(), query)
	return res
}

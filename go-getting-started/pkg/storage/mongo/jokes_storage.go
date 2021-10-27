package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type JokesStorage struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

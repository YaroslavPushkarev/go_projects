package storage

import (
	"context"

	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type JokesInterface interface {
	FindId(string) (models.Joke, error)
	GetJokes(interface{}) ([]models.Joke, error)
	InsertJoke(models.Joke) (*mongo.InsertOneResult, error)
}

type JokesHandler struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

func (j *JokesHandler) InsertJoke(joke models.Joke) (*mongo.InsertOneResult, error) {
	res, err := j.Collection.InsertOne(j.Ctx, joke)

	if err != nil {
		return res, err
	}

	return res, err
}

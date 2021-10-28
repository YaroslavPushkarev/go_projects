package mongo

import (
	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (j JokesStorage) GetJokes(filter interface{}) ([]models.Joke, error) {
	jokes := []models.Joke{}

	if filter == nil {
		filter = bson.M{}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(5)).SetSkip(int64(5))

	cursor, err := j.Collection.Find(j.Ctx, filter, findOptions)
	if err != nil {
		return jokes, err
	}

	for cursor.Next(j.Ctx) {
		row := models.Joke{}
		err = cursor.Decode(&row)
		if err != nil {
			panic(err)
		}
		jokes = append(jokes, row)
	}

	return jokes, nil
}

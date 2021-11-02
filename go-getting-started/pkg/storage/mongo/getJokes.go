package mongo

import (
	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (j JokesStorage) GetJokes(filter interface{}, limit, skip int) ([]models.Joke, error) {
	jokes := []models.Joke{}

	if filter == nil {
		filter = bson.D{}
	}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit)).SetSkip(int64(skip))

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

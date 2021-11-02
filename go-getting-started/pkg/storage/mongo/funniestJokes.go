package mongo

import (
	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (j JokesStorage) FunniestJokes(filter interface{}, limit int) ([]models.Joke, error) {

	jokes := []models.Joke{}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: -1}}).SetLimit(int64(limit))

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

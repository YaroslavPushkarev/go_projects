package mongo

import (
	"log"

	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (j JokesStorage) FunniestJokes(interface{}) ([]models.Joke, error) {

	jokes := []models.Joke{}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: -1}}).SetLimit(int64(20))

	cursor, err := j.Collection.Find(j.Ctx, findOptions)
	if err != nil {
		return jokes, err
	}

	if err = cursor.All(j.Ctx, &jokes); err != nil {
		log.Fatal(err)
	}

	return jokes, nil
}

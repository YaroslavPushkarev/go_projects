package mongo

import (
	"context"

	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (j JokesStorage) FindJoke(s string) ([]models.Joke, error) {
	jokes := []models.Joke{}

	query := bson.M{
		"title": bson.M{
			"$regex": primitive.Regex{
				Pattern: s,
				Options: "i",
			},
		},
	}

	cursor, err := j.Collection.Find(j.Ctx, query)
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var joke models.Joke
		err := cursor.Decode(&joke)
		if err != nil {
			panic(err)
		}
		jokes = append(jokes, joke)
	}

	return jokes, nil
}

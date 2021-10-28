package mongo

import (
	"fmt"

	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (j JokesStorage) Random(filter interface{}) ([]models.Joke, error) {

	jokes := []models.Joke{}

	pipeline := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 10}}}}

	cursor, err := j.Collection.Aggregate(j.Ctx, mongo.Pipeline{pipeline})
	if err != nil {
		fmt.Println(err)
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

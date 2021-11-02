package mongo

import (
	"github.com/heroku/go-getting-started/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (j JokesStorage) FindId(id string) (models.Joke, error) {
	joke := models.Joke{}

	err := j.Collection.FindOne(j.Ctx, bson.M{"id": id}).Decode(&joke)
	if err != nil {
		return joke, err
	}

	return joke, nil
}

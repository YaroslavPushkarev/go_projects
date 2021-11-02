package mongo

import (
	"fmt"

	"github.com/heroku/go-getting-started/pkg/models"
)

func (j JokesStorage) InsertJoke(joke models.Joke) error {

	res, err := j.Collection.InsertOne(j.Ctx, joke)

	fmt.Println(res)

	if err != nil {
		return err
	}

	return err
}

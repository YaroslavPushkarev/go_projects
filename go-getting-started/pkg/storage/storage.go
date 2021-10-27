package storage

import (
	"github.com/heroku/go-getting-started/pkg/models"
)

type JokesStorage interface {
	FindId(string) (models.Joke, error)
	GetJokes(interface{}) ([]models.Joke, error)
	InsertJoke(models.Joke) error
}

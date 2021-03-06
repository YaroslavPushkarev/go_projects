package storage

import (
	"github.com/heroku/go-getting-started/pkg/models"
)

type JokesStorage interface {
	FindId(string) (models.Joke, error)
	FindJoke(string) ([]models.Joke, error)
	GetJokes(interface{}, int, int) ([]models.Joke, error)
	Random(interface{}, int) ([]models.Joke, error)
	FunniestJokes(interface{}, int) ([]models.Joke, error)
	InsertJoke(models.Joke) error
}

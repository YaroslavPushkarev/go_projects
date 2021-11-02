package api

import "github.com/heroku/go-getting-started/pkg/storage"

type JokesHandler struct {
	Storage storage.JokesStorage
}

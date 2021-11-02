package mongo_test

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/heroku/go-getting-started/pkg/api"
	"github.com/heroku/go-getting-started/pkg/models"
	"github.com/heroku/go-getting-started/pkg/storage/mongo"
	"github.com/stretchr/testify/assert"
)

func TestControllers_InsertJoke_CheckingAddedJokes(t *testing.T) {
	tt := []struct {
		name string
		want models.Joke
		id   string
	}{
		{
			name: "joke 1",
			want: models.Joke{
				Body:  "asdada",
				ID:    "sdfadsa",
				Score: 4,
				Title: "asd",
			},
			id: "sdfadsa",
		},
		{
			name: "empty title",
			want: models.Joke{
				Body:  "asdaa",
				ID:    "gsdfsf",
				Score: 4,
				Title: "",
			},
			id: "gsdfsf",
		},
		{
			name: "empty body",
			want: models.Joke{
				Body:  "",
				ID:    "rdsf",
				Score: 4,
				Title: "asd",
			},
			id: "rdsf",
		},
	}

	collection := mongo.ConnectDB("mongodb://localhost:27017")

	str := mongo.JokesStorage{
		Collection: collection,
	}

	client := api.JokesHandler{
		Storage: str,
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			err := client.Storage.InsertJoke(tc.want)
			assert.Nil(t, err)

			joke, err := client.Storage.FindId(tc.id)

			assert.Nil(t, err)
			assert.Equal(t, tc.want, joke)
		})
	}
}

func TestControllers_InsertJoke_WhatIfInsertIdenticalID(t *testing.T) {

	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	randomID := fmt.Sprintf("%X", b)

	tt := []struct {
		name   string
		want   models.Joke
		insert models.Joke
		id     string
	}{
		{
			name: "identical id",
			insert: models.Joke{
				Body:  "A Sunday school teacher is concerned that his students might be a litt",
				ID:    randomID,
				Score: 4,
				Title: "his hand",
			},
		},
	}
	collection := mongo.ConnectDB("mongodb://localhost:27017")

	str := mongo.JokesStorage{
		Collection: collection,
	}

	client := api.JokesHandler{
		Storage: str,
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			err := client.Storage.InsertJoke(tc.insert)
			assert.Nil(t, err)

			err = client.Storage.InsertJoke(tc.insert)
			assert.Error(t, err)

		})
	}
}

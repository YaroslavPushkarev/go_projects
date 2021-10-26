package storage_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/heroku/go-getting-started/pkg/models"
	storage "github.com/heroku/go-getting-started/pkg/storage/mongo"
	"github.com/stretchr/testify/assert"
)

func TestControllers_InsertJoke(t *testing.T) {
	tt := []struct {
		name string
		want models.Joke
		id   string
	}{
		{
			name: "Big number",
			want: models.Joke{
				Body:  "asdada",
				ID:    "sdfadsa",
				Score: 4,
				Title: "asd",
			},
			id: "sdfadsa",
		},
		{
			name: "zero",
			want: models.Joke{
				Body:  "asdaa",
				ID:    "gsdfsf",
				Score: 4,
				Title: "asd",
			},
			id: "gsdfsf",
		},
		{
			name: "Empty title",
			want: models.Joke{
				Body:  "asdada",
				ID:    "rdsf",
				Score: 4,
				Title: "asd",
			},
			id: "rdsf",
		},
	}

	ctx := context.TODO()
	collection := storage.ConnectDB("mongodb://localhost:27017")

	client := &storage.JokesHandler{
		Collection: collection,
		Ctx:        ctx,
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			_, err := client.InsertJoke(tc.want)
			assert.Nil(t, err)

			joke, err := client.FindId(tc.id)

			assert.Nil(t, err)
			assert.Equal(t, tc.want, joke)
		})
	}
}

func TestControllers_InsertIdenticalID(t *testing.T) {

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
			name: randomID,
			insert: models.Joke{
				Body:  "A Sunday school teacher is concerned that his students might be a litt",
				ID:    "dg",
				Score: 4,
				Title: "his hand",
			},
		},
	}
	ctx := context.TODO()
	collection := storage.ConnectDB("mongodb://localhost:27017")

	client := &storage.JokesHandler{
		Collection: collection,
		Ctx:        ctx,
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			_, err := client.InsertJoke(tc.insert)
			assert.Nil(t, err)

			_, err = client.InsertJoke(tc.insert)
			assert.Error(t, err)

		})
	}
}

package controllers

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/heroku/go-getting-started/config"
	"github.com/heroku/go-getting-started/models"
	"github.com/stretchr/testify/assert"
)

func TestControllers_InsertData(t *testing.T) {
	tt := []struct {
		name string
		want models.Joke
	}{
		{
			name: "Big number",
			want: models.Joke{
				Body:  "asdada",
				ID:    "sdfadsa",
				Score: 4,
				Title: "asd",
			},
		},
		{
			name: "zero",
			want: models.Joke{
				Body:  "asdaa",
				ID:    "gsdfsf",
				Score: 4,
				Title: "asd",
			},
		},
		{
			name: "Empty title",
			want: models.Joke{
				Body:  "asdada",
				ID:    "rdsf",
				Score: 4,
				Title: "asd",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var joke models.Joke

			collection := config.ConnectDB("mongodb://localhost:27017")

			res, err := InsertData(collection, tc.want)
			assert.Nil(t, err)

			err = FindId(collection, map[string]interface{}{"_id": res.InsertedID}).Decode(&joke)

			assert.Nil(t, err)
			assert.Equal(t, tc.want, joke)
		})
	}
}

func TestControllers_InsertDataID(t *testing.T) {

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
			name: "id",
			insert: models.Joke{
				Body:  "A Sunday school teacher is concerned that his students might be a litt",
				ID:    randomID,
				Score: 4,
				Title: "his hand",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			collection := config.ConnectDB("mongodb://localhost:27017")

			_, err := InsertData(collection, tc.insert)
			assert.Nil(t, err)

			_, err = InsertData(collection, tc.insert)
			assert.Error(t, err)

		})
	}
}

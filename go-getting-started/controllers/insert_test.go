package controllers

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestControllers_InsertData(t *testing.T) {
	tt := []struct {
		name  string
		title string
	}{
		{
			name:  "title string",
			title: "3",
		},
		{
			name:  "empty title",
			title: "",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			clientOptions := options.Client().ApplyURI("mongodb+srv://jokesdb:jokesdb@joke.kxki9.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx, clientOptions)

			if err != nil {
				log.Fatal(err)
			}

			collection := client.Database("Jokes").Collection("jokes")
			res, err := InsertData(collection, Joke{"asdada", tc.title, 3, "4324"})

			assert.Nil(t, err)
			assert.IsType(t, &mongo.InsertOneResult{}, res)
		})
	}
}

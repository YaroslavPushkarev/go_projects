package controllers

import (
	"context"
	"testing"

	"github.com/heroku/go-getting-started/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mockCollection struct {
}

func (m *mockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	c := &mongo.InsertOneResult{}
	return c, nil
}

func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	c := &mongo.SingleResult{}
	return c
}

func TestControllers_InsertData(t *testing.T) {
	tt := []struct {
		name  string
		score int
		title string
	}{
		{
			name:  "Big number",
			score: 23424324243,
			title: "a",
		},
		{
			name:  "less than zero",
			score: -32,
			title: "b",
		},
		{
			name:  "zero",
			score: 0,
			title: "c",
		},
		{
			name:  "Empty title",
			title: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockCollection := &mockCollection{}

			res, err := InsertData(mockCollection, models.Joke{Body: "asdada", ID: "sdf", Score: tc.score, Title: tc.title})

			assert.Nil(t, err)
			assert.IsType(t, &mongo.InsertOneResult{}, res)
		})
	}
}

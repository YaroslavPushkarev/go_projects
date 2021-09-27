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

func TestControllers_InsertData(t *testing.T) {
	mockCollection := &mockCollection{}
	res, err := InsertData(mockCollection, models.Joke{Body: "asdada", ID: "sdf", Score: 3, Title: "4324"})

	assert.Nil(t, err)
	assert.IsType(t, &mongo.InsertOneResult{}, res)

}

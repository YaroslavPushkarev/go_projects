package controllers

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *mockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	c := &mongo.SingleResult{}
	return c
}

// func TestControllers_FindID(t *testing.T) {s
// 	mockCollection := &mockCollection{}
// 	query := map[string]interface{}{"id": "5tz2wj"}
// 	red := FindId(mockCollection, query)

// 	assert.Equal(t, res, got)

// }

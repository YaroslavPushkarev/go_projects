package api

import "go.mongodb.org/mongo-driver/mongo"

type Joke struct {
	Body  string `json:"body" bson:"body"`
	ID    string `json:"id" bson:"id"`
	Score int    `json:"score" bson:"score"`
	Title string `json:"title" bson:"title"`
}

type Pagination struct {
	Skip  int
	Limit int
}

type JokesHandler struct {
	jokes      []Joke
	collection *mongo.Collection
}

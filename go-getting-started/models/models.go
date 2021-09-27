package models

type Joke struct {
	Body  string `json:"body" bson:"body"`
	ID    string `json:"id" bson:"id"`
	Score int    `json:"score" bson:"score"`
	Title string `json:"title" bson:"title"`
}

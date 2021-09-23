// package controllers

// import (
// 	"context"
// 	"log"
// 	"testing"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func TestInsertJoke(t *testing.T) {
// 	storage := []Joke{
// 		{
// 			Body:  "Now I have to say \"Leroy can you please paint the fence?\"",
// 			ID:    "5tz52q",
// 			Score: 1,
// 			Title: "I hate how you cant even say black paint anymore",
// 		},
// 		{
// 			Body:  "Pizza doesn't scream when you put it in the oven .\n\nI'm so sorry.",
// 			ID:    "5tz4dd",
// 			Score: 0,
// 			Title: "What's the difference between a Jew in Nazi Germany and pizza ?",
// 		},
// 	}
// 	clientOptions := options.Client().ApplyURI("mongodb+srv://jokesdb:jokesdb@joke.kxki9.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, clientOptions)
// 	collection := client.Database("Jokes").Collection("jokes")

// 	res := CreateJoke

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

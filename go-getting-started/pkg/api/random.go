func (j jokesHandler) randomJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jokes []models.Joke

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		fmt.Println(err)
	}

	pipeline := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: pagination.Limit}}}}

	cursor, err := j.collection.Aggregate(context.TODO(), mongo.Pipeline{pipeline})
	if err != nil {
		fmt.Println(err)
	}

	if err = cursor.All(context.TODO(), &jokes); err != nil {
		panic(err)
	}

	if err = json.NewEncoder(w).Encode(jokes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

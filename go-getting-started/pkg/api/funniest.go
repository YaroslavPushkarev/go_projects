func (j jokesHandler) funniest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var jokes []models.Joke

	pagination, err := j.parseSkipAndLimit(w, r)
	if err != nil {
		fmt.Println(err)
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: -1}}).SetLimit(int64(pagination.Limit))

	cursor, err := j.collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = cursor.All(context.TODO(), &jokes); err != nil {
		log.Fatal(err)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func (j JokesHandler) search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	search := r.URL.Query().Get("search")

	var jokes []models.Joke

	query := bson.M{
		"title": bson.M{
			"$regex": primitive.Regex{
				Pattern: search,
				Options: "i",
			},
		},
	}

	cursor, err := j.collection.Find(context.TODO(), query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for cursor.Next(context.TODO()) {
		var joke models.Joke
		err := cursor.Decode(&joke)
		if err != nil {
			panic(err)
		}
		jokes = append(jokes, joke)
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

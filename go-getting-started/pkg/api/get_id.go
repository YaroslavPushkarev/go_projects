func (j JokesHandler) getId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	_, err := j.collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		fmt.Println(err)
	}

	var jokes models.Joke

	id := r.URL.Query().Get("id")

	query := map[string]interface{}{"id": id}

	err = controllers.FindId(j.collection, query).Decode(&jokes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(jokes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

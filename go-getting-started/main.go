package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Pagination struct {
	Skip  int
	Limit int
}
type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Score int    `json:"score"`
	Body  string `json:"body"`
}

type jokesHandler struct {
	jokes []Joke
}

const maxLimit = 1

func (j jokesHandler) parseSkipAndLimit(w http.ResponseWriter, r *http.Request) (Pagination, error) {
	leng := len(j.jokes)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 1
	}
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		skip = 1
	}
	if skip > leng {
		w.WriteHeader(http.StatusNoContent)
		return Pagination{}, nil
	}
	if skip < 0 {
		skip = 1
	}

	if limit > leng {
		limit = maxLimit
	}

	if limit < 0 {
		limit = maxLimit
	}

	pagination := Pagination{Skip: skip, Limit: limit}
	return pagination, nil
}

func (p jokesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(p.jokes) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(p.jokes)
		return
	}

	if len(p.jokes) == 1 {
		json.NewEncoder(w).Encode(p.jokes)
		return
	}

	jokes := p.jokes

	pagination, err := p.parseSkipAndLimit(w, r)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]

	json.NewEncoder(w).Encode(res)

}

// func getJoke(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	id := r.URL.Query().Get("id")

// 	for _, item := range jokes {
// 		if item.ID == id {
// 			json.NewEncoder(w).Encode(item)
// 		}
// 	}
// }

// func randomJokes(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	rand.Seed(time.Now().UTC().UnixNano())

// 	rand.Shuffle(len(jokes), func(i, j int) { jokes[i], jokes[j] = jokes[j], jokes[i] })
// 	pagination, err := parseSkipAndLimit(r)
// 	fmt.Fprintln(w, err)
// 	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]

// 	json.NewEncoder(w).Encode(res)
// }

// func funniest(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	sort.SliceStable(jokes, func(i, j int) bool {
// 		return jokes[i].Score > jokes[j].Score
// 	})

// 	pagination, err := parseSkipAndLimit(r)
// 	fmt.Fprintln(w, err)
// 	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]

// 	json.NewEncoder(w).Encode(res)
// }

// func search(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var arr []string

// 	value := r.URL.Query().Get("input")

// 	for _, item := range jokes {
// 		if strings.Contains(item.Title, value) {
// 			arr = append(arr, item.Title)
// 		}
// 	}

// 	json.NewEncoder(w).Encode(arr)

// }

func main() {
	content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := []Joke{}
	json.Unmarshal(content, &jokes)

	mux := http.NewServeMux()
	mux.Handle("/jokes", jokesHandler{jokes})
	log.Fatal(http.ListenAndServe(":8080", mux))
	// http.HandleFunc("/", index)
	// http.HandlerFunc("/jokes", pizzasHandler{&data})
	// http.HandleFunc("/jokes/", getJoke)
	// http.HandleFunc("/jokes/random/", randomJokes)
	// http.HandleFunc("/jokes/funniest/", funniest)
	// http.HandleFunc("/search", search)
	// http.ListenAndServe(":8000", nil)
}

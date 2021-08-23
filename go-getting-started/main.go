package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

type Pagination struct {
	Skip  int
	Limit int
}

func parseSkipAndLimit(r *http.Request) (Pagination, error) {

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		err = fmt.Errorf("err")
		return Pagination{}, err
	}

	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil || skip < 1 {
		err = fmt.Errorf("err")
		return Pagination{}, err
	}

	pagination := Pagination{Skip: skip, Limit: limit}

	return pagination, nil
}

type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Score int    `json:"score"`
	// Body  string `json:"body"`
}

type Jokes []Joke

type jokesHandler struct {
	jokes *Jokes
}

func (p jokesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jokes := *p.jokes
	pagination, err := parseSkipAndLimit(r)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	if pagination.Skip <= 0 && pagination.Limit <= 0 {
		res := jokes[1:5]
		json.NewEncoder(w).Encode(res)
	}

	if pagination.Skip > len(jokes) && pagination.Limit > len(jokes) {
		res := jokes[1:5]
		json.NewEncoder(w).Encode(res)
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
	// content, _ := ioutil.ReadFile("reddit_jokes.json")
	jokes := Jokes{
		Joke{
			ID:    "1",
			Title: "Why women need legs?",
			Score: 12,
		},
		Joke{
			ID:    "2",
			Title: "I recently went to America....",
			Score: 11,
		},
		Joke{
			ID:    "3",
			Title: "I recently went to America....",
			Score: 10,
		},
		Joke{
			ID:    "5",
			Title: "I recently went to America....",
			Score: 4,
		},
		Joke{
			ID:    "6",
			Title: "I recently went to America....",
			Score: 16,
		},
		Joke{
			ID:    "7",
			Title: "I recently went to America....",
			Score: 34,
		},
		Joke{
			ID:    "8",
			Title: "I recently went to America....",
			Score: 5,
		},
	}
	// json.Unmarshal(content, &jokes)

	mux := http.NewServeMux()
	mux.Handle("/jokes", jokesHandler{&jokes})
	log.Fatal(http.ListenAndServe(":8080", mux))
	// http.HandleFunc("/", index)
	// http.HandlerFunc("/jokes", pizzasHandler{&data})
	// http.HandleFunc("/jokes/", getJoke)
	// http.HandleFunc("/jokes/random/", randomJokes)
	// http.HandleFunc("/jokes/funniest/", funniest)
	// http.HandleFunc("/search", search)
	// http.ListenAndServe(":8000", nil)
}

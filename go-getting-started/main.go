package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/heroku/x/hmetrics/onload"
)

type Joke struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Score int    `json:"score"`
	Body  string `json:"body"`
}

var jokes []Joke

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
		return Pagination{Limit: 2}, err
	}

	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil || skip < 1 {
		err = fmt.Errorf("err")
		return Pagination{}, err
	}

	pagination := Pagination{Skip: skip, Limit: limit}

	return pagination, nil
}

func getJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pagination, err := parseSkipAndLimit(r)

	// s := r.URL.Query().Get("skip")
	// if s == "" {
	// 	http.Error(w, "missing value", http.StatusBadRequest)
	// 	return
	// }

	// skip, err := strconv.Atoi(s)
	// if err != nil {
	// 	http.Error(w, "not a number: "+s, http.StatusBadRequest)
	// 	return
	// }

	// l := r.URL.Query().Get("limit")
	// if err != nil  {
	// 	    http.Error(w, err.Error(), http.StatusBadRequest)
	// 	    return
	// }

	// limit, err := strconv.Atoi(l)
	// 	if err != nil {
	// 	http.Error(w, "not a number: "+l, http.StatusBadRequest)

	// return
	// }

	fmt.Fprintln(w, err)
	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]
	json.NewEncoder(w).Encode(res)
}

func getJoke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")

	for _, item := range jokes {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func randomJokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rand.Seed(time.Now().UTC().UnixNano())

	rand.Shuffle(len(jokes), func(i, j int) { jokes[i], jokes[j] = jokes[j], jokes[i] })
	pagination, err := parseSkipAndLimit(r)
	fmt.Fprintln(w, err)
	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]

	json.NewEncoder(w).Encode(res)
}

func funniest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sort.SliceStable(jokes, func(i, j int) bool {
		return jokes[i].Score > jokes[j].Score
	})

	pagination, err := parseSkipAndLimit(r)
	fmt.Fprintln(w, err)
	res := jokes[pagination.Skip : pagination.Limit+pagination.Skip]

	json.NewEncoder(w).Encode(res)
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var arr []string

	value := r.URL.Query().Get("input")

	for _, item := range jokes {
		if strings.Contains(item.Title, value) {
			arr = append(arr, item.Title)
		}
	}

	json.NewEncoder(w).Encode(arr)

}

func main() {

	content, _ := ioutil.ReadFile("reddit_jokes.json")
	json.Unmarshal(content, &jokes)

	http.HandleFunc("/", index)
	http.HandleFunc("/jokes", getJokes)
	http.HandleFunc("/jokes/", getJoke)
	http.HandleFunc("/jokes/random/", randomJokes)
	http.HandleFunc("/jokes/funniest/", funniest)
	http.HandleFunc("/search", search)

	http.ListenAndServe(":8000", nil)
}

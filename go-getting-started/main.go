package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	// "os"
	"sort"
	"strings"
	"text/template"
	"time"
	"github.com/bmizerany/pat"


	// "github.com/gin-gonic/gin"

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

func getJokes(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", jokes)
}

func getJoke(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	for _, item := range jokes {
		if item.ID == id {
			tpl.ExecuteTemplate(w, "id.html", item)
			return
		}
	}
}

func randomJokes(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())

	rand.Shuffle(len(jokes), func(i, j int) { jokes[i], jokes[j] = jokes[j], jokes[i] })

	tpl.ExecuteTemplate(w, "random.html", jokes)
}

func fun(w http.ResponseWriter, r *http.Request) {
	sort.SliceStable(jokes, func(i, j int) bool {
		return jokes[i].Score > jokes[j].Score
	})
	tpl.ExecuteTemplate(w, "fun.html", jokes)
}


func search(w http.ResponseWriter, r *http.Request) {
	var arr []string

	value := r.FormValue("q")

	for _, item := range jokes {
		contain := strings.Contains(item.Title, value)
			if contain == true {
			arr = append(arr, item.Title)		
		}			
	}
	tpl.ExecuteTemplate(w, "search.html", arr)
}

func main() {

	content, _ := ioutil.ReadFile("reddit_jokes.json")
	json.Unmarshal(content, &jokes)

	//port := os.Getenv("PORT")

	//if port == "" {
    //		log.Fatal(":8000")
	//}

	r := pat.New()

	http.Handle("/", r)
	r.Get("/", http.HandlerFunc(index))
	r.Get("/jokes", http.HandlerFunc(getJokes))
	r.Get("/jokes/:id", http.HandlerFunc(getJoke))
	r.Get("/jokes/random/", http.HandlerFunc(randomJokes))
	r.Get("/jokes/funniest/", http.HandlerFunc(fun))
	r.Get("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	r.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.Post("/search", http.HandlerFunc(search))
	//log.Fatal(http.ListenAndServe(":"+port, r))
	err := http.ListenAndServe(":8000", nil)
	 if err != nil {
	 	log.Fatal("ListenAndServe: ", err)
	 }
}

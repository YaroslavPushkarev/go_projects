package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJokesHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      []Joke
		want       string
		statusCode int
	}{
		{
			name:       "without jokes",
			method:     http.MethodGet,
			input:      []Joke{},
			want:       "[]",
			statusCode: http.StatusNoContent,
		},
		{
			name:   "with jokes",
			method: http.MethodGet,
			input: []Joke{
				{
					ID:    "1",
					Title: "Foo",
					Score: 10,
					Body:  "sdfsf",
				},
			},
			want:       `[{"id":"1","title":"Foo","score":10,"body":"sdfsf"}]`,
			statusCode: http.StatusOK,
		},
		{
			name:       "with bad method",
			method:     http.MethodPost,
			input:      []Joke{},
			want:       "Method not allowed",
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/jokes?skip=1&limit=3", nil)
			responseRecorder := httptest.NewRecorder()

			jokesHandler{tc.input}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}

func TestJokesHandler_pagination(t *testing.T) {

	storage := []Joke{
		{
			Body:  "Now I have to say \"Leroy can you please paint the fence?\"",
			ID:    "5tz52q",
			Score: 1,
			Title: "I hate how you cant even say black paint anymore",
		},
		{
			Body:  "Pizza doesn't scream when you put it in the oven .\n\nI'm so sorry.",
			ID:    "5tz4dd",
			Score: 0,
			Title: "What's the difference between a Jew in Nazi Germany and pizza ?",
		},
		{
			Body:  "...and being there really helped me learn about American culture. So I visited a shop and as I was leaving, the Shopkeeper said \"Have a nice day!\" But I didn't so I sued him.",
			ID:    "5tz319",
			Score: 0,
			Title: "I recently went to America....",
		},
		{
			Body:  "A Sunday school teacher is concerned that his students might be a little confused about Jesus, so he asks his class, \u201cWhere is Jesus today?\u201d\nBrian raises his hand and says, \u201cHe\u2019s in Heaven.\u201d\n\nSusan answers, \u201cHe\u2019s in my heart.\u201d\n\nLittle Johnny waves his hand furiously and blurts out, \u201cHe\u2019s in our bathroom!\u201d\n\nThe teacher is surprised by this answer and asks Little Johnny how he knows this.\n\n\u201cWell,\u201d Little Johnny says, \u201cevery morning, my Dad gets up, bangs on the bathroom door and yells \u2018Jesus Christ, are you still in there?'\u201d",
			ID:    "5tz2wj",
			Score: 1,
			Title: "Brian raises his hand and says, \u201cHe\u2019s in Heaven.\u201d",
		},
		{
			Body:  "He got caught trying to sell the two books to a freshman.",
			ID:    "5tz1pc",
			Score: 0,
			Title: "You hear about the University book store worker who was charged for stealing $20,000 worth of books?",
		},
	}
	tt := []struct {
		name       string
		want       []Joke
		statusCode int
		skip       int
		limit      int
	}{
		{
			name:       "skip 0 limit 1 - expect id 5tz52q",
			want:       []Joke{storage[0]},
			statusCode: http.StatusOK,
			skip:       0,
			limit:      1,
		},
		{
			name:       "skip 100000000000000 limit 1",
			want:       []Joke{storage[4]},
			statusCode: http.StatusOK,
			skip:       100000000000,
			limit:      1,
		},
		{
			name: "skip -2 limit 4",
			want: []Joke{
				{
					Body:  "Pizza doesn't scream when you put it in the oven .\n\nI'm so sorry.",
					ID:    "5tz4dd",
					Score: 0,
					Title: "What's the difference between a Jew in Nazi Germany and pizza ?",
				},
				{
					Body:  "...and being there really helped me learn about American culture. So I visited a shop and as I was leaving, the Shopkeeper said \"Have a nice day!\" But I didn't so I sued him.",
					ID:    "5tz319",
					Score: 0,
					Title: "I recently went to America....",
				},
				{
					Body:  "A Sunday school teacher is concerned that his students might be a little confused about Jesus, so he asks his class, \u201cWhere is Jesus today?\u201d\nBrian raises his hand and says, \u201cHe\u2019s in Heaven.\u201d\n\nSusan answers, \u201cHe\u2019s in my heart.\u201d\n\nLittle Johnny waves his hand furiously and blurts out, \u201cHe\u2019s in our bathroom!\u201d\n\nThe teacher is surprised by this answer and asks Little Johnny how he knows this.\n\n\u201cWell,\u201d Little Johnny says, \u201cevery morning, my Dad gets up, bangs on the bathroom door and yells \u2018Jesus Christ, are you still in there?'\u201d",
					ID:    "5tz2wj",
					Score: 1,
					Title: "Brian raises his hand and says, \u201cHe\u2019s in Heaven.\u201d",
				},
				{
					Body:  "He got caught trying to sell the two books to a freshman.",
					ID:    "5tz1pc",
					Score: 0,
					Title: "You hear about the University book store worker who was charged for stealing $20,000 worth of books?",
				},
			},
			statusCode: http.StatusOK,
			skip:       -2,
			limit:      4,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			uri := fmt.Sprintf("/jokes?skip=%d&limit=%d", tc.skip, tc.limit)

			request := httptest.NewRequest(http.MethodGet, uri, nil)
			responseRecorder := httptest.NewRecorder()

			jokesHandler{storage}.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			b, err := ioutil.ReadAll(responseRecorder.Body)
			if err != nil {
				t.Errorf("reading reponse body: %v", err)
			}

			got := []Joke{}
			err = json.Unmarshal(b, &got)

			if err != nil {
				t.Errorf("reading reponse body: %v", err)
			}
			assert.Equal(t, tc.want, got)

		})
	}
}

// package main

// import (
// 	"bytes"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"
// 	"testing"
// )

// func TestJokes(t *testing.T) {
// 	tt := []struct {
// 		name  string
// 		value string
// 		err   string
// 		input int
// 	}{
// 		{name: "get value", value: "5", input: 5},
// 		{name: "missing value", value: "", err: "missing value"},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			req, err := http.NewRequest("GET", "localhost:8000/jokes?skip="+tc.value, nil)
// 			if err != nil {
// 				t.Fatalf("could not send GET request: %v", err)
// 			}

// 			rec := httptest.NewRecorder()

// 			getJokes(rec, req)

// 			res := rec.Result()

// 			defer res.Body.Close()

// 			b, err := ioutil.ReadAll(res.Body)
// 			if err != nil {
// 				t.Fatalf("could not response: %v", err)
// 			}

// 			if tc.err != "" {
// 				if res.StatusCode != http.StatusBadRequest {
// 					t.Errorf("expected status Bad Request; got %v", res.StatusCode)
// 				}

// 				if msg := string(bytes.TrimSpace(b)); msg != tc.err {
// 					t.Errorf("expected message %q; got %q", tc.err, msg)
// 				}
// 				return
// 			}

// 			if res.StatusCode != http.StatusOK {
// 				t.Errorf("expected status OK; got %v", res.Status)
// 			}

// 			d, err := strconv.Atoi(string(bytes.TrimSpace(b)))
// 			if err != nil {
// 				t.Fatalf("expected an integer; got %s", b)
// 			}

// 			if d != tc.input {
// 				t.Fatalf("expected %v; got %v", tc.input, d)
// 			}
// 		})
// 	}
// }

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJokesHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		input      *Jokes
		want       string
		statusCode int
	}{
		{
			name:       "without jokes",
			method:     http.MethodGet,
			input:      &Jokes{},
			want:       "Error: No jokes found",
			statusCode: http.StatusNotFound,
		},
		{
			name:   "with jokes",
			method: http.MethodGet,
			input: &Jokes{
				Joke{
					ID:    "1",
					Title: "Foo",
					Score: 10,
				},
			},
			want:       `[{"id":"1", "title":"Foo","score":10}]`,
			statusCode: http.StatusOK,
		},
		{
			name:       "with bad method",
			method:     http.MethodPost,
			input:      &Jokes{},
			want:       "Method not allowed",
			statusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/jokes", nil)
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

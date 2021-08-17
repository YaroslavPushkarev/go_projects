package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"strconv"
	"io/ioutil"
)

func TestJokes(t *testing.T){
	tt := []struct {
		name string
		value string
		status int

	}{
		{name: "get value", value:"5", status: http.StatusOK},
	}

	for _, tc := range tt {
		req, err := http.NewRequest("GET", "localhost:8000/jokes?skip="+tc.value, nil)
		if err != nil {
			t.Fatalf("could not send GET request: %v", err)
		}

		rec := httptest.NewRecorder()

		getJokes(rec, req)

		res := rec.Result()		

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status OK; got %v", res.Status)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("could not response: %v", err)
		}

		d, err := strconv.Atoi(string(bytes.TrimSpace(b)))
		if err != nil {
			t.Fatalf("expected an integer; got %s", b)
		}

		if d != 5  {
			t.Fatalf("expected 5; got %v", d)
		}
	}
}
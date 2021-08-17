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
	req, err := http.NewRequest("GET", "localhost:8000/jokes?skip=5", nil)
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
		t.Fatalf("expected double to be 5; got %v", d)
	}
}
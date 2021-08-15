package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)



func TestTokenProcessing(t *testing.T) {
    rr := httptest.NewRecorder()
    r, err := http.NewRequest("GET", "http://golang.org/", nil)
    if err != nil {
        t.Fatal(err)
    }

    r.URL.Query().Add("name", "300")

    handler := http.HandlerFunc(getJokes)
    handler.ServeHTTP(rr, r)

    if code := rr.Code; code != http.StatusOK {
        t.Fatalf("handler did not return correct status: want %v got %v",
            http.StatusOK, code)
    }

}
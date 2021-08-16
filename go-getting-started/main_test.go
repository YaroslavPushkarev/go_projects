package main

import (
    "fmt"
    "net/http"
	"net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)



func TestJokes(t *testing.T) {
    testCases := []struct  {
        skip string
        
    }{
        {
            skip: "5",
         
        },
        {
            skip: "20",
        
        },
    }

    handler := http.HandlerFunc(getJokes)

    for _, tc := range testCases{
        t.Run(tc.skip, func(t *testing.T){
            record := httptest.NewRecorder()
            request, _ := http.NewRequest("GET", fmt.Sprintf("/jokes?skip=%d", tc.skip), nil)
            handler.ServeHTTP(record, request)
            assert.Equal(t, tc.want, record.Body.Bytes())
        })
    }
    // rr := httptest.NewRecorder()
    // r, err := http.NewRequest("GET", "http://golang.org/", nil)
    // if err != nil {
    //     t.Fatal(err)
    // }

    // r.URL.Query().Add("skip", "5")
	// r.URL.Query().Add("limit", "20")

    // handler := http.HandlerFunc(getJokes)
    // handler.ServeHTTP(rr, r)

    // if code := rr.Code; code != http.StatusOK {
    //     t.Fatalf("handler did not return correct status: want %v got %v",
    //         http.StatusOK, code)
    // }

}
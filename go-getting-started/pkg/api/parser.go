package api

import (
	"net/http"
	"strconv"
)

func (j JokesHandler) parseSkipAndLimit(w http.ResponseWriter, r *http.Request) (Pagination, error) {
	leng := len(j.jokes)

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		limit = 1
	}
	skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		skip = 1
	}
	if skip > leng {
		w.WriteHeader(http.StatusBadRequest)
		return Pagination{}, nil
	}
	if skip < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return Pagination{}, nil
	}

	if limit > leng {
		limit = leng - skip
	}

	if limit < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return Pagination{}, nil
	}

	pagination := Pagination{Skip: skip, Limit: limit}
	return pagination, nil
}

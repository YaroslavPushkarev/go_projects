package api

import (
	"net/http"
	"strconv"

	"github.com/heroku/go-getting-started/pkg/models"
)

func ParseSkipAndLimit(w http.ResponseWriter, r *http.Request) (models.Pagination, error) {
	leng := len("10000")

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
		return models.Pagination{}, nil
	}
	if skip < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return models.Pagination{}, nil
	}

	if limit > leng {
		limit = leng - skip
	}

	if limit < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return models.Pagination{}, nil
	}

	pagination := models.Pagination{Skip: skip, Limit: limit}
	return pagination, nil
}

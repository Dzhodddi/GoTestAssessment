package store

import (
	"net/http"
	"strconv"
)

type PaginatedQuery struct {
	Limit  int `json:"limit" validate:"gte=1,lte=100"`
	Offset int `json:"offset" validate:"gte=0"`
}

func (fq PaginatedQuery) Parse(r *http.Request) (PaginatedQuery, error) {
	q := r.URL.Query()

	limit := q.Get("limit")
	offset := q.Get("offset")
	if limit != "" {
		o, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = o
	}

	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}

	return fq, nil
}

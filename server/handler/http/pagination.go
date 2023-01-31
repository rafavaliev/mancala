package http

import (
	"context"
	"net/http"
	"strconv"
)

const (
	defaultLimit  = 10
	defaultOffset = 0
	maxLimit      = 100
	minLimit      = 1

	// PaginationCtxKey refers to the context key that stores the pagination
	PaginationCtxKey string = "pagination"
)

type Pagination struct {
	Limit  int
	Offset int
}

// paginationMiddleware is used to extract offset and limit from the url query
func paginationMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			limitQ := r.URL.Query().Get("limit")
			offsetQ := r.URL.Query().Get("offset")

			limit := defaultLimit
			offset := defaultOffset
			var err error

			if limitQ != "" {
				limit, err = strconv.Atoi(limitQ)
				if err != nil {
					limit = defaultLimit
				}
				if limit > maxLimit {
					limit = maxLimit
				}
				if limit < minLimit {
					limit = defaultLimit
				}
			}

			if offsetQ != "" {
				offset, err = strconv.Atoi(offsetQ)
				if err != nil {
					offset = defaultOffset
				}
				if offset < 0 {
					offset = defaultOffset
				}
			}

			ctx := context.WithValue(r.Context(), PaginationCtxKey, Pagination{
				Limit:  limit,
				Offset: offset,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

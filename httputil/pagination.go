package httputil

import (
	"net/url"
	"strconv"
)

// Pagination represents the pagination values.
// It allows you to paginate your database queries using a cursor based pagination.
type Pagination struct {
	PerPage   int
	Cursor    int
	Next      bool
	FirstPage bool
}

// GetPagination returns the pagination values from the query parameters for use in your database query.
func GetPagination(q url.Values) Pagination {
	var (
		pp, _     = strconv.Atoi(q.Get("perPage"))
		cursor, _ = strconv.Atoi(q.Get("cursor"))
		next, _   = strconv.ParseBool(q.Get("next"))
		firstPage = false
	)

	if pp > 100 {
		pp = 100
	}

	if !next && cursor < 1 {
		firstPage = true
	}

	return Pagination{
		PerPage:   pp,
		Cursor:    cursor,
		Next:      next,
		FirstPage: firstPage,
	}
}

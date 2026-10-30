package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/bookindex"
)

const (
	defaultSearchHitsPerPage = 20
	maxSearchHitsPerPage     = 100
)

// searchBooks handles GET /api/books.
//
// Query parameters:
//
//	q                                   full-text query
//	genre, publisher, format, tag       facet values (repeat to select many)
//	yearMin, yearMax, priceMin, priceMax, ratingMin
//	sort                                relevance (default), newest, oldest, price_asc, price_desc, rating, popular
//	page (from 1), hitsPerPage (max 100)
//
//declscope:shared // main.go routes to it
func (h *handler) searchBooks(w http.ResponseWriter, r *http.Request) {
	params, err := parseSearchParams(r.URL.Query())
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := bookindex.Search(r.Context(), h.client, params)
	if err != nil {
		log.Printf("search: %v", err)
		writeErrorResponse(w, http.StatusBadGateway, "search failed")
		return
	}
	writeJSONResponse(w, http.StatusOK, res)
}

func parseSearchParams(q url.Values) (bookindex.SearchParams, error) {
	p := bookindex.SearchParams{
		Query:       q.Get("q"),
		Genres:      q["genre"],
		Publishers:  q["publisher"],
		Formats:     q["format"],
		Tags:        q["tag"],
		Page:        1,
		HitsPerPage: defaultSearchHitsPerPage,
	}

	var ok bool
	if p.Sort, ok = bookindex.ParseSortOrder(q.Get("sort")); !ok {
		return p, fmt.Errorf("unknown sort %q", q.Get("sort"))
	}

	var errs []error
	intParam := func(name string) *int {
		s := q.Get(name)
		if s == "" {
			return nil
		}
		v, err := strconv.Atoi(s)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s must be an integer", name))
			return nil
		}
		return &v
	}
	p.YearMin = intParam("yearMin")
	p.YearMax = intParam("yearMax")
	p.PriceMin = intParam("priceMin")
	p.PriceMax = intParam("priceMax")
	if page := intParam("page"); page != nil {
		p.Page = int64(max(*page, 1))
	}
	if n := intParam("hitsPerPage"); n != nil {
		p.HitsPerPage = int64(min(max(*n, 1), maxSearchHitsPerPage))
	}
	if s := q.Get("ratingMin"); s != "" {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			errs = append(errs, errors.New("ratingMin must be a number"))
		} else {
			p.RatingMin = &v
		}
	}
	return p, errors.Join(errs...)
}

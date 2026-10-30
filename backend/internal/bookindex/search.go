package bookindex

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/meilisearch/meilisearch-go"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/book"
)

// SearchHighlightPreTag and SearchHighlightPostTag wrap matched words in SearchHit.Formatted.
const (
	searchHighlightPreTag  = "<mark>"
	searchHighlightPostTag = "</mark>"
)

// SearchParams is the input of Search.
// Values inside one facet field are OR-ed; different fields are AND-ed.
type SearchParams struct {
	Query       string
	Genres      []string
	Publishers  []string
	Formats     []string
	Tags        []string
	YearMin     *int
	YearMax     *int
	PriceMin    *int
	PriceMax    *int
	RatingMin   *float64
	Sort        SortOrder
	Page        int64
	HitsPerPage int64
}

// SearchHit is a book with search metadata.
type SearchHit struct {
	book.Book
	Formatted    SearchHighlights `json:"_formatted"`
	RankingScore float64          `json:"_rankingScore"`
}

// SearchHighlights holds copies of the searchable fields with matches highlighted.
type SearchHighlights struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

// SearchFacetStat is the min and max of a numeric field over the matched books.
type SearchFacetStat struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// SearchResult is the output of Search.
type SearchResult struct {
	Hits             []SearchHit                 `json:"hits"`
	Query            string                      `json:"query"`
	Page             int64                       `json:"page"`
	HitsPerPage      int64                       `json:"hitsPerPage"`
	TotalHits        int64                       `json:"totalHits"`
	TotalPages       int64                       `json:"totalPages"`
	ProcessingTimeMs int64                       `json:"processingTimeMs"`
	Facets           map[string]map[string]int64 `json:"facets"`
	FacetStats       map[string]SearchFacetStat  `json:"facetStats"`
}

// Search runs the book search.
//
// It sends one multi-search request. The first query returns the hits.
// For each facet that has selected values, one more query counts that facet
// without its own filter. This is "disjunctive faceting": after you check
// "SF", the genre list still shows how many books the other genres have.
func Search(ctx context.Context, client meilisearch.ServiceManager, p SearchParams) (*SearchResult, error) {
	// Numeric attributes are requested as facets only to get their min and max.
	statAttributes := []string{"publishedYear", "price"}
	facets := append([]string{}, statAttributes...)
	for _, ff := range p.facetFilters() {
		facets = append(facets, ff.attribute)
	}

	hitsPerPage := p.HitsPerPage
	queries := []*meilisearch.SearchRequest{{
		IndexUID:              indexUID,
		Query:                 p.Query,
		Filter:                p.filter(""),
		Sort:                  p.Sort.meilisearchSort(),
		Page:                  p.Page,
		HitsPerPage:           &hitsPerPage,
		Facets:                facets,
		AttributesToHighlight: []string{"title", "author", "description"},
		HighlightPreTag:       searchHighlightPreTag,
		HighlightPostTag:      searchHighlightPostTag,
		ShowRankingScore:      true,
	}}
	var disjunctive []string
	for _, ff := range p.facetFilters() {
		if len(ff.values) == 0 {
			continue
		}
		zero := int64(0)
		queries = append(queries, &meilisearch.SearchRequest{
			IndexUID:    indexUID,
			Query:       p.Query,
			Filter:      p.filter(ff.attribute),
			HitsPerPage: &zero,
			Facets:      []string{ff.attribute},
		})
		disjunctive = append(disjunctive, ff.attribute)
	}

	res, err := client.MultiSearchWithContext(ctx, &meilisearch.MultiSearchRequest{Queries: queries})
	if err != nil {
		return nil, err
	}

	main := res.Results[0]
	out := &SearchResult{
		Hits:             make([]SearchHit, 0, len(main.Hits)),
		Query:            p.Query,
		Page:             main.Page,
		HitsPerPage:      main.HitsPerPage,
		TotalHits:        main.TotalHits,
		TotalPages:       main.TotalPages,
		ProcessingTimeMs: main.ProcessingTimeMs,
	}
	if err := main.Hits.DecodeInto(&out.Hits); err != nil {
		return nil, fmt.Errorf("decode hits: %w", err)
	}
	if err := json.Unmarshal(main.FacetDistribution, &out.Facets); err != nil {
		return nil, fmt.Errorf("decode facets: %w", err)
	}
	for _, attr := range statAttributes {
		delete(out.Facets, attr)
	}
	if len(main.FacetStats) > 0 {
		if err := json.Unmarshal(main.FacetStats, &out.FacetStats); err != nil {
			return nil, fmt.Errorf("decode facet stats: %w", err)
		}
	}
	for i, attr := range disjunctive {
		r := res.Results[i+1]
		var dist map[string]map[string]int64
		if err := json.Unmarshal(r.FacetDistribution, &dist); err != nil {
			return nil, fmt.Errorf("decode facets for %s: %w", attr, err)
		}
		out.Facets[attr] = dist[attr]
		out.ProcessingTimeMs = max(out.ProcessingTimeMs, r.ProcessingTimeMs)
	}
	return out, nil
}

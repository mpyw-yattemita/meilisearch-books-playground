package bookindex

// SortOrder is the sort option exposed by the API.
type SortOrder string

// Supported sort orders.
const (
	SortRelevance SortOrder = "relevance"
	SortNewest    SortOrder = "newest"
	SortOldest    SortOrder = "oldest"
	SortPriceAsc  SortOrder = "price_asc"
	SortPriceDesc SortOrder = "price_desc"
	SortRating    SortOrder = "rating"
	SortPopular   SortOrder = "popular"
)

var sortRules = map[SortOrder][]string{
	SortRelevance: nil,
	SortNewest:    {"publishedYear:desc"},
	SortOldest:    {"publishedYear:asc"},
	SortPriceAsc:  {"price:asc"},
	SortPriceDesc: {"price:desc"},
	SortRating:    {"rating:desc", "reviewCount:desc"},
	SortPopular:   {"reviewCount:desc"},
}

// ParseSortOrder returns the SortOrder for s, or false if s is unknown.
// An empty string means SortRelevance.
func ParseSortOrder(s string) (SortOrder, bool) {
	if s == "" {
		return SortRelevance, true
	}
	_, ok := sortRules[SortOrder(s)]
	return SortOrder(s), ok
}

// meilisearchSort returns the Meilisearch "sort" parameter for o.
// SortRelevance returns nil, which keeps the ranking rules' order.
//
//declscope:shared // search.go builds the request from it
func (o SortOrder) meilisearchSort() []string {
	return sortRules[o]
}

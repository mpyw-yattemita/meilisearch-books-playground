package bookindex

import (
	"fmt"
	"strconv"
	"strings"
)

// facetFilter is one multi-select facet in the UI and the values selected in it.
//
//declscope:shared // search.go asks one facet query per selected filter
type facetFilter struct {
	attribute string
	values    []string
}

//declscope:shared // search.go asks one facet query per selected filter
func (p *SearchParams) facetFilters() []facetFilter {
	return []facetFilter{
		{"genre", p.Genres},
		{"publisher", p.Publishers},
		{"format", p.Formats},
		{"tags", p.Tags},
	}
}

// filter builds a Meilisearch filter: the outer array is AND, inner arrays are OR.
// The facet named by exclude is left out (see Search for why).
//
//declscope:shared // search.go puts it in every query
func (p *SearchParams) filter(exclude string) []any {
	var f []any
	for _, ff := range p.facetFilters() {
		if ff.attribute == exclude || len(ff.values) == 0 {
			continue
		}
		or := make([]string, 0, len(ff.values))
		for _, v := range ff.values {
			or = append(or, fmt.Sprintf("%s = %s", ff.attribute, quoteFilterValue(v)))
		}
		f = append(f, or)
	}
	if p.YearMin != nil {
		f = append(f, fmt.Sprintf("publishedYear >= %d", *p.YearMin))
	}
	if p.YearMax != nil {
		f = append(f, fmt.Sprintf("publishedYear <= %d", *p.YearMax))
	}
	if p.PriceMin != nil {
		f = append(f, fmt.Sprintf("price >= %d", *p.PriceMin))
	}
	if p.PriceMax != nil {
		f = append(f, fmt.Sprintf("price <= %d", *p.PriceMax))
	}
	if p.RatingMin != nil {
		f = append(f, "rating >= "+strconv.FormatFloat(*p.RatingMin, 'f', -1, 64))
	}
	if len(f) == 0 {
		return nil
	}
	return f
}

func quoteFilterValue(s string) string {
	return `"` + strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(s) + `"`
}

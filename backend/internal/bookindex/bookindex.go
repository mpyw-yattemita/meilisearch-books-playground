//declscope:core
//declscope:shared // the index definition every file works against

// Package bookindex owns the "books" Meilisearch index: its settings,
// bulk loading, and the search query used by the HTTP API.
package bookindex

import (
	"github.com/meilisearch/meilisearch-go"
)

// indexUID is the Meilisearch index name.
const indexUID = "books"

// primaryKey is the document field Meilisearch uses as the ID.
const primaryKey = "id"

// settings returns the index settings applied by Seed.
func settings() *meilisearch.Settings {
	return &meilisearch.Settings{
		// Order matters: matches in title rank higher than in description.
		SearchableAttributes: []string{"title", "author", "tags", "publisher", "description", "isbn"},
		FilterableAttributes: []string{"genre", "publisher", "format", "tags", "publishedYear", "price", "rating", "pages"},
		SortableAttributes:   []string{"publishedYear", "price", "rating", "reviewCount"},
		// Without this, short CJK text such as "京都" may be tokenized as Chinese.
		LocalizedAttributes: []*meilisearch.LocalizedAttributes{
			{Locales: []string{"jpn"}, AttributePatterns: []string{"title", "author", "publisher", "description", "tags"}},
		},
		Synonyms: map[string][]string{
			"golang":   {"Go"},
			"ts":       {"TypeScript"},
			"k8s":      {"Kubernetes"},
			"postgres": {"PostgreSQL"},
		},
		Pagination: &meilisearch.Pagination{MaxTotalHits: 100_000},
		Faceting:   &meilisearch.Faceting{MaxValuesPerFacet: 100},
	}
}

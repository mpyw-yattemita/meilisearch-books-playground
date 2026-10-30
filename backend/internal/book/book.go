// Package book defines the book document stored in Meilisearch.
package book

// Book is one document in the "books" index.
type Book struct {
	ID            int      `json:"id"`
	ISBN          string   `json:"isbn"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Publisher     string   `json:"publisher"`
	Genre         string   `json:"genre"`
	Format        string   `json:"format"`
	Tags          []string `json:"tags"`
	Description   string   `json:"description"`
	PublishedYear int      `json:"publishedYear"`
	Pages         int      `json:"pages"`
	Price         int      `json:"price"`
	Rating        float64  `json:"rating"`
	ReviewCount   int      `json:"reviewCount"`
}

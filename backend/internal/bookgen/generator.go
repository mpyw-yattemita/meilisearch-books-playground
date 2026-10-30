// Package bookgen generates dummy books. The same seed always yields the same books.
package bookgen

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/book"
)

// Generator builds dummy books. The same seed always yields the same books.
type Generator struct {
	rng *rand.Rand
}

// NewGenerator returns a Generator seeded with seed.
func NewGenerator(seed uint64) *Generator {
	return &Generator{rng: rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))}
}

// GenerateBooks returns count books with IDs 1..count.
func (g *Generator) GenerateBooks(count int) []book.Book {
	books := make([]book.Book, 0, count)
	for i := range count {
		books = append(books, g.generateBook(i+1))
	}
	return books
}

// generateBook returns one book with the given ID.
func (g *Generator) generateBook(id int) book.Book {
	genre := pick(g.rng, genres)
	year := 1975 + g.weightedRecent(52)
	pages := 80 + g.rng.IntN(520)
	if genre.name == "絵本" {
		pages = 24 + g.rng.IntN(40)
	}
	return book.Book{
		ID:            id,
		ISBN:          g.isbn(),
		Title:         g.title(genre),
		Author:        pick(g.rng, familyNames) + " " + pick(g.rng, givenNames),
		Publisher:     pick(g.rng, publishers),
		Genre:         genre.name,
		Format:        pick(g.rng, formats),
		Tags:          g.tags(genre),
		Description:   g.description(genre),
		PublishedYear: year,
		Pages:         pages,
		Price:         g.price(pages),
		Rating:        math.Round((1.5+g.rng.Float64()*3.5)*10) / 10,
		ReviewCount:   int(math.Pow(g.rng.Float64(), 3) * 5000),
	}
}

// weightedRecent returns 0..span-1, biased toward larger values (recent years).
func (g *Generator) weightedRecent(span int) int {
	return int(math.Sqrt(g.rng.Float64()) * float64(span))
}

func (g *Generator) isbn() string {
	digits := make([]int, 0, 13)
	digits = append(digits, 9, 7, 8, 4) // 978-4: Japanese ISBN prefix
	for range 8 {
		digits = append(digits, g.rng.IntN(10))
	}
	sum := 0
	for i, d := range digits {
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}
	digits = append(digits, (10-sum%10)%10)
	var sb strings.Builder
	for _, d := range digits {
		sb.WriteByte(byte('0' + d))
	}
	return sb.String()
}

func (g *Generator) title(genre genreSpec) string {
	pattern := pick(g.rng, genre.titlePatterns)
	r := strings.NewReplacer(
		"{adj}", pick(g.rng, adjectives),
		"{noun}", pick(g.rng, genre.nouns),
		"{noun2}", pick(g.rng, genre.nouns),
		"{place}", pick(g.rng, places),
		"{num}", fmt.Sprint(2+g.rng.IntN(99)),
	)
	return r.Replace(pattern)
}

func (g *Generator) tags(genre genreSpec) []string {
	n := 1 + g.rng.IntN(3)
	seen := make(map[string]struct{}, n)
	tags := make([]string, 0, n)
	for range n {
		t := pick(g.rng, genre.tags)
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		tags = append(tags, t)
	}
	return tags
}

func (g *Generator) description(genre genreSpec) string {
	r := strings.NewReplacer(
		"{noun}", pick(g.rng, genre.nouns),
		"{place}", pick(g.rng, places),
		"{adj}", pick(g.rng, adjectives),
	)
	return r.Replace(pick(g.rng, genre.descriptions)) + pick(g.rng, closingSentences)
}

func (g *Generator) price(pages int) int {
	base := 500 + pages*5 + g.rng.IntN(1500)
	return base / 10 * 10
}

func pick[T any](rng *rand.Rand, s []T) T {
	return s[rng.IntN(len(s))]
}

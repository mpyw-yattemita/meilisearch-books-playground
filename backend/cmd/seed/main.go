// Command seed generates dummy books and loads them into Meilisearch.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/bookgen"
	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/bookindex"
	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/config"
)

func main() {
	count := flag.Int("count", 50_000, "number of books to generate")
	seed := flag.Uint64("seed", 42, "random seed (same seed, same books)")
	batchSize := flag.Int("batch", 10_000, "documents per add-documents request")
	dump := flag.Bool("dump", false, "print books as JSON to stdout instead of loading them")
	flag.Parse()

	books := bookgen.NewGenerator(*seed).GenerateBooks(*count)
	if *dump {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(books); err != nil {
			log.Fatal(err)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	start := time.Now()
	log.Printf("loading %d books (batch size %d)...", len(books), *batchSize)
	if err := bookindex.Seed(ctx, config.NewMeilisearchClient(), books, *batchSize); err != nil {
		log.Fatal(err)
	}
	log.Printf("done in %s", time.Since(start).Round(time.Millisecond))
}

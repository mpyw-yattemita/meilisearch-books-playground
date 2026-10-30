//declscope:core

// Command server serves the book search API on top of Meilisearch.
package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/meilisearch/meilisearch-go"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/config"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	h := &handler{client: config.NewMeilisearchClient()}
	srv := &http.Server{
		Addr:              *addr,
		Handler:           logRequests(h.routes()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("listening on %s", *addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// handler holds what every endpoint needs. Each endpoint lives in its own file.
//
//declscope:shared // every endpoint file is a method set on it
type handler struct {
	client meilisearch.ServiceManager
}

func (h *handler) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/books", h.searchBooks)
	mux.HandleFunc("GET /api/books/{id}", h.getBook)
	mux.HandleFunc("GET /api/health", h.health)
	return mux
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.RequestURI(), time.Since(start).Round(time.Microsecond))
	})
}

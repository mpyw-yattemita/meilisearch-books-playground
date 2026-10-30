package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/meilisearch/meilisearch-go"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/bookindex"
)

// getBook handles GET /api/books/{id}.
//
//declscope:shared // main.go routes to it
func (h *handler) getBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	b, err := bookindex.GetBook(r.Context(), h.client, id)
	if meiliErr, ok := errors.AsType[*meilisearch.Error](err); ok && meiliErr.StatusCode == http.StatusNotFound {
		writeErrorResponse(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		log.Printf("get book: %v", err)
		writeErrorResponse(w, http.StatusBadGateway, "get book failed")
		return
	}
	writeJSONResponse(w, http.StatusOK, b)
}

package bookindex

import (
	"context"
	"strconv"

	"github.com/meilisearch/meilisearch-go"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/book"
)

// GetBook returns the book with the given ID.
func GetBook(ctx context.Context, client meilisearch.ServiceManager, id int) (*book.Book, error) {
	var b book.Book
	if err := client.Index(indexUID).GetDocumentWithContext(ctx, strconv.Itoa(id), nil, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

package bookindex

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"

	"github.com/mpyw-yattemita/meilisearch-books-playground/backend/internal/book"
)

// Seed drops the index, recreates it with settings, and loads books in batches.
func Seed(ctx context.Context, client meilisearch.ServiceManager, books []book.Book, batchSize int) error {
	// Deleting an index that does not exist yet fails with index_not_found; that is fine.
	if err := waitTask(ctx, client, func() (*meilisearch.TaskInfo, error) {
		return client.DeleteIndexWithContext(ctx, indexUID)
	}, "index_not_found"); err != nil {
		return fmt.Errorf("delete index: %w", err)
	}
	if err := waitTask(ctx, client, func() (*meilisearch.TaskInfo, error) {
		return client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{Uid: indexUID, PrimaryKey: primaryKey})
	}); err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	index := client.Index(indexUID)
	// Apply settings before adding documents so Meilisearch indexes only once.
	if err := waitTask(ctx, client, func() (*meilisearch.TaskInfo, error) {
		return index.UpdateSettingsWithContext(ctx, settings())
	}); err != nil {
		return fmt.Errorf("update settings: %w", err)
	}

	pk := primaryKey
	tasks, err := index.AddDocumentsInBatchesWithContext(ctx, books, batchSize, &meilisearch.DocumentOptions{PrimaryKey: &pk})
	if err != nil {
		return fmt.Errorf("add documents: %w", err)
	}
	for _, t := range tasks {
		if err := waitTask(ctx, client, func() (*meilisearch.TaskInfo, error) { return &t, nil }); err != nil {
			return fmt.Errorf("add documents (task %d): %w", t.TaskUID, err)
		}
	}
	return nil
}

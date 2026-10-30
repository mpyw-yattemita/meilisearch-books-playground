package bookindex

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

const taskPollInterval = 200 * time.Millisecond

// waitTask enqueues a task and blocks until Meilisearch finishes it.
// A failure whose code is in ignoredCodes counts as success.
//
//declscope:shared // seed.go runs every index change through it
func waitTask(ctx context.Context, client meilisearch.ServiceManager, enqueue func() (*meilisearch.TaskInfo, error), ignoredCodes ...string) error {
	info, err := enqueue()
	if err != nil {
		return err
	}
	task, err := client.WaitForTaskWithContext(ctx, info.TaskUID, taskPollInterval)
	if err != nil {
		return err
	}
	if task.Status == meilisearch.TaskStatusSucceeded || slices.Contains(ignoredCodes, task.Error.Code) {
		return nil
	}
	return errors.New(task.Error.Code + ": " + task.Error.Message)
}

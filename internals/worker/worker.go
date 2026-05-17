package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/imsumedhaa/go-taskqueue/internals/database"
)

type Worker struct {
	ID   int
	Repo *database.TaskRepository
}

func NewWorker(id int, repo *database.TaskRepository) *Worker {
	return &Worker{
		ID:   id,
		Repo: repo,
	}
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("Worker %d started\n", w.ID)

	for {
		id, taskType, payload, retryCount, maxRetry, err := w.Repo.GetPendingTask(ctx)

		if err != nil {
			fmt.Printf("Worker %d: no pending task\n", w.ID)
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf(
			"Worker %d picked task %s (%s) [retry %d/%d]\n",
			w.ID, id, taskType, retryCount, maxRetry,
		)

		err = w.executeTask(taskType, payload)

		if err != nil {
			fmt.Printf(
				"Worker %d failed task %s: %v\n",
				w.ID, id, err,
			)

			updateErr := w.Repo.UpdateTaskStatus(ctx, id, "failed", err.Error())

			if updateErr != nil {
				fmt.Printf(
					"Worker %d failed to update task status to 'failed' %s: %v\n",
					w.ID, id, updateErr,
				)
			}

			// retry logic -> worker only retry the task if it doesn't cross the max retries limit

			if retryCount+1 < maxRetry {

				fmt.Printf(
					"Worker %d retrying task %s (%d/%d)\n",
					w.ID,
					id,
					retryCount+1,
					maxRetry,
				)

				time.Sleep(3 * time.Second)

				retryErr := w.Repo.UpdateTaskStatus(ctx, id, "pending", "")

				if retryErr != nil {
					fmt.Printf(
						"Worker %d failed to update task status to 'pending' %s: %v\n",
						w.ID, id, retryErr,
					)

				}
			} else {
				fmt.Printf(
					"Worker %d: task %s reached max retries (%d). Marked permanently failed\n",
					w.ID,
					id,
					maxRetry,
				)
			}
			continue // here continue means if it reached the max retries, skip the task and move forward to the next

		}

		err = w.Repo.UpdateTaskStatus(ctx, id, "completed", "")
		if err != nil {
			fmt.Printf("Worker %d failed to update status to 'completed' for task %s: %v\n",
				w.ID, id, err)
		}

		fmt.Printf("Worker %d successfully completed task %s (%s)\n", w.ID, id, taskType)
	}
}

func (w *Worker) executeTask(tasktype string, payload []byte) error {
	switch tasktype {

	case "send_email":
		fmt.Printf("Workner %d is sending email..\n", w.ID)
		time.Sleep(2 * time.Second)
		return nil

	case "generate_report":
		fmt.Printf("Workner %d is generating report..\n", w.ID)
		time.Sleep(2 * time.Second)
		return nil

	default:
		return fmt.Errorf("Worker %d received unknown task type %s\n", w.ID)

	}

}

func StartWorkerPool(ctx context.Context, numWorkers int, repo *database.TaskRepository) {

	for i := 1; i <= numWorkers; i++ {
		workers := NewWorker(i, *&repo)
		go workers.Start(ctx)
	}

}

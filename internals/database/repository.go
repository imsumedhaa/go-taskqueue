package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) InsertTask(ctx context.Context, tasktype string, payload interface{}) (string, error) {

	id := uuid.New().String()

	payloadByte, err := json.Marshal(payload)

	if err != nil {
		return "", fmt.Errorf("Error in marshal the payload %w\n", err)
	}

	query := `insert into tasks(id, type, payload) values ($1, $2, $3)`

	_, err = r.db.Pool.Exec(ctx, query, id, tasktype, payloadByte)

	if err != nil {
		return "", fmt.Errorf("Error in inserting %w\n", err)
	}

	return id, nil
}

func (r *TaskRepository) GetPendingTask(ctx context.Context) (string, string, []byte, int, int, error) {

	var id string
	var tasktype string
	var payload []byte
	var retryCount int
	var maxRetry int

	tx, err := r.db.Pool.Begin(ctx)

	if err != nil {
		return "", "", nil, 0, 0, fmt.Errorf("failed to begin the transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id, type, payload, retry_count, max_retries
		FROM tasks
		WHERE status = 'pending'
		ORDER BY created_at
		FOR UPDATE SKIP LOCKED
		LIMIT 1
	`

	err = tx.QueryRow(ctx, query).Scan(&id, &tasktype, &payload, &retryCount, &maxRetry)
	if err != nil {
		return "", "", nil, 0, 0, fmt.Errorf("failed to fetch the pendind task: %w", err)
	}

	updateQuery := `
	update tasks set status = 'processing', updated_at = NOW() where id = $1
	`

	_, err = tx.Exec(ctx, updateQuery, id)

	if err != nil {
		return "", "", nil, 0, 0, fmt.Errorf("failed to update the status of the id: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", nil, 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return id, tasktype, payload, retryCount, maxRetry, nil

}

// Update task status to completed / failed / retry

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, id string, status string, errMsg string) error {

	// if the status is failed then retry count will increase by 1, otherwise same as before

	query := ` 
	UPDATE tasks 
	SET status = $1, 
	last_error = $2, 
	retry_count = CASE 
	  WHEN $1 = 'failed' THEN retry_count + 1 
	     ELSE retry_count 
		 END,
	updated_at = NOW()
	WHERE id = $3 `

	_, err := r.db.Pool.Exec(ctx, query, status, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

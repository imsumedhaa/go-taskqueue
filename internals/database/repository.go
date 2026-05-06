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

	payloadbyte, err := json.Marshal(payload)

	if err != nil {
		return "", fmt.Errorf("Error in marshal the payload %w\n", err)
	}

	query := `insert into tasks(id, tasktype, payload) values ($1, $2, $3)`

	_, err = r.db.Pool.Exec(ctx, query, id, tasktype, payloadbyte)

	if err != nil {
		return "", fmt.Errorf("Error in inserting %w\n", err)
	}

	return id, nil
}

func (r *TaskRepository) GetPendingTask(ctx context.Context) (string, string, interface{}, error) {

	var id string
	var tasktype string
	var payload []byte

	tx, err := r.db.Pool.Begin(ctx)

	if err != nil {
		return "", "", nil, fmt.Errorf("failed to begin the transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		select id, type, payload from tasks where status='pending' order by for update skip locked  created_at limit 1
	`

	err = tx.QueryRow(ctx, query).Scan(&id, &tasktype, &payload)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to fetch the pendind task: %w", err)
	}

	updateQuery := `
	update tasks set status = 'processing', updated_at = NOW() where id = $1
	`

	_, err = tx.Exec(ctx, updateQuery, id)

	if err != nil {
		return "", "", nil, fmt.Errorf("failed to update the status of the id: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return id, tasktype, payload, nil

}

// Update task status to completed / failed / retry

func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, id string, status string, errMsg string) error {

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

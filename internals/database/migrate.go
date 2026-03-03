package database

import (
	"context"
	"fmt"
)

func (db *DB) CreateTaskTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks(
	id UUID primary key,
	type text not null,
	payload jsonb not null,
	status text not null default 'pending',
	retry_count int default 0,
	max_retries int default 5,
	last_error text, 
	created_at timestamp not null default now(),
	updated_at timestamp not null default now()
	)
	`
	_, err := db.Pool.Exec(ctx, query)

	if err != nil {
		return fmt.Errorf("Failed to create tasks table: %w\n", err)
	}

	fmt.Println("tasks table created..")
	return nil
}

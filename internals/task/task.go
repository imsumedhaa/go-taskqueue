package task

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        string
	Type      string
	Payload   []byte
	CreatedAt time.Time
}

func NewTask(taskType string, payload []byte) Task {
	return Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Payload:   payload,
		CreatedAt: time.Now(),
	}
}

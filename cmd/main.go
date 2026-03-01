package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/imsumedhaa/go-taskqueue/internals/queue"
	"github.com/imsumedhaa/go-taskqueue/internals/task"
	"github.com/imsumedhaa/go-taskqueue/internals/worker"
)

func main() {

	fmt.Println("Starting Task Queue System...")

	// Create Queue
	q := queue.NewQueue(10)

	// 2. Start Worker Pool --> worker 3
	worker.StartWorkerPool(3, q)

	// 3. Create Tasks
	emailPayload, _ := json.Marshal(map[string]string{
		"to":      "user@gmail.com",
		"subject": "Hello",
	})

	reportPayload, _ := json.Marshal(map[string]string{
		"report": "monthly_sales",
	})

	task1 := task.NewTask("send_email", emailPayload)
	task2 := task.NewTask("generate_report", reportPayload)

	q.Enqueue(task1)
	q.Enqueue(task2)

	fmt.Println("Tasks submitted successfully")

	time.Sleep(10 * time.Second)

}

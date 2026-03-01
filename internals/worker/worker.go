package worker

import (
	"fmt"
	"time"

	"github.com/imsumedhaa/go-taskqueue/internals/queue"
	"github.com/imsumedhaa/go-taskqueue/internals/task"
)

type Worker struct {
	ID    int
	Queue queue.Queue
}

func NewWorker(id int, q queue.Queue) *Worker {
	return &Worker{
		ID:    id,
		Queue: q,
	}
}

func (w *Worker) start() {
	fmt.Printf("Worker %d started \n", w.ID)

	for t := range w.Queue.Tasks() {
		fmt.Printf("Worker %d picked task %s (%s)\n", w.ID, t.ID, t.Type)

		w.execute(t)
	}
}

func (w *Worker) execute(t task.Task) {
	switch t.Type {

	case "send_email":
		fmt.Printf("Workner %d is sending email..\n", w.ID)
		time.Sleep(2 * time.Second)

	case "generate_report":
		fmt.Printf("Workner %d is generating report..\n", w.ID)
		time.Sleep(2 * time.Second)

	default:
		fmt.Printf("Worker %d received unknown task type %s\n", w.ID, t.Type)
	}

	fmt.Printf("Worker %d finished task %s\n", w.ID, t.ID)
}

func StartWorkerPool(numWorkers int, q *queue.Queue) {

	for i := 1; i <= numWorkers; i++ {
		workers := NewWorker(i, *q)
		go workers.start()
	}

}

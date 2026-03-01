package queue

import "github.com/imsumedhaa/go-taskqueue/internals/task"

type Queue struct {
	tasks chan task.Task
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{
		tasks: make(chan task.Task, bufferSize),
	}
}

func (q *Queue)Enqueue(t task.Task){
	q.tasks <- t
}

func (q *Queue)Tasks() <- chan task.Task{    //<- chan task.Task is the return type --> read-only channel
	return q.tasks
}
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/imsumedhaa/go-taskqueue/internals/database"
	"github.com/imsumedhaa/go-taskqueue/internals/worker"
)

func main() {

	fmt.Println("Starting Task Queue System...")

	ctx := context.Background()

	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	err = db.CreateTaskTable(ctx)
	if err != nil {
		panic(err)
	}

	repo := database.NewTaskRepository(db)

	_, err = repo.InsertTask(ctx, "send-email", map[string]string{
		"to":      "sumedha@gmail.com",
		"subject": "hello world",
	})
	if err != nil {
		fmt.Println("Error creating email task:", err)
	}

	repo.InsertTask(ctx, "generate_report", map[string]string{
		"report": "monthly_salary",
	})
	if err != nil {
		fmt.Println("Error creating report task:", err)
	}

	fmt.Println("Tasks inserted into db")

	worker.StartWorkerPool(ctx, 3, repo)

	time.Sleep(20 * time.Second)

}

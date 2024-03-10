package main

import (
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynq/x/rate"
	"go-labs/middleware/asyncq/task"
	"log"
)

// workers.go
func main() {
	redisClientOpt := asynq.RedisClientOpt{Addr: "localhost:6379"}
	srv := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			Concurrency: 1,
			IsFailure: func(err error) bool {
				return !task.IsRateLimitError(err)
			},
			RetryDelayFunc: task.RetryDelay,
		},
	)

	sema := rate.NewSemaphore(redisClientOpt, "reminder-email-semaphore", 1)
	defer func() {
		_ = sema.Close()
	}()
	processor := task.NewEmailProcessor(sema)
	ReminderEmailHandlers := asynq.NewServeMux()
	ReminderEmailHandlers.Use(processor.ReminderEmailMiddleware)
	ReminderEmailHandlers.HandleFunc(task.TypeReminderEmail, processor.HandleReminderEmailTask)

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TypeWelcomeEmail, processor.HandleWelcomeEmailTask)
	mux.Handle(task.TypeReminderEmail, ReminderEmailHandlers)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}

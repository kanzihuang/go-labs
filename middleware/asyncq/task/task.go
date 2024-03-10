package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq/x/rate"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeWelcomeEmail  = "email:welcome"
	TypeReminderEmail = "email:reminder"
)

type RateLimitError struct {
	RetryIn time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limited (retry in  %v)", e.RetryIn)
}

func IsRateLimitError(err error) bool {
	var rateLimitError *RateLimitError
	return errors.As(err, &rateLimitError)

}

func RetryDelay(n int, err error, task *asynq.Task) time.Duration {
	var rateLimitErr *RateLimitError
	if errors.As(err, &rateLimitErr) {
		return rateLimitErr.RetryIn
	}
	return asynq.DefaultRetryDelayFunc(n, err, task)
}

// Task payload for any email related tasks.
type emailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

func NewWelcomeEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeWelcomeEmail, payload), nil
}

func NewReminderEmailTask(id int) (*asynq.Task, error) {
	payload, err := json.Marshal(emailTaskPayload{UserID: id})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeReminderEmail, payload), nil
}

func NewEmailProcessor(sema *rate.Semaphore) *EmailProcessor {
	return &EmailProcessor{
		semaphore: sema,
	}
}

type EmailProcessor struct {
	semaphore *rate.Semaphore
}

func (p *EmailProcessor) HandleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	time.Sleep(time.Second * 5)
	log.Printf(" [*] Send Welcome Email to User %d", payload.UserID)
	return nil
}

func (p *EmailProcessor) HandleReminderEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload emailTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	time.Sleep(time.Second * 5)
	log.Printf(" [*] Send Reminder Email to User %d", payload.UserID)
	return nil
}

func (p *EmailProcessor) ReminderEmailMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		ok, err := p.semaphore.Acquire(ctx)
		if err != nil {
			return err
		}
		if !ok {
			return &RateLimitError{RetryIn: time.Second * 3}
		}
		// Make sure to release the token once we're done.
		defer func() {
			// todo: 如果 Release 执行失败，必须自动清理过期 token，不然会导致 token 被长期占用
			_ = p.semaphore.Release(ctx)
		}()

		return h.ProcessTask(ctx, task)
	})
}

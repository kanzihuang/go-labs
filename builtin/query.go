package builtin

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Queue[T any] interface {
	Enqueue(ctx context.Context, val T) error
	Dequeue() (T, error)
}

var _ Queue[int] = &SliceQueue[int]{}

const querySize = 256

func NewSliceQueue[T any](size int) Queue[T] {
	mutex := &sync.Mutex{}
	return &SliceQueue[T]{
		buffer:    make([]T, size),
		size:      size,
		mutex:     mutex,
		condFull:  sync.NewCond(mutex),
		condEmpty: sync.NewCond(mutex),
	}
}

type SliceQueue[T any] struct {
	buffer    []T
	r         int
	w         int
	size      int
	zero      T
	mutex     *sync.Mutex
	condFull  *sync.Cond
	condEmpty *sync.Cond
}

var (
	errTimeout = errors.New("SliceQueue: enqueue timeout")
)

func (s *SliceQueue[T]) Enqueue(ctx context.Context, val T) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	done := false
	go func() {
		select {
		case <-ctx.Done():
			done = true
			s.condFull.Broadcast()
		}
	}()
	for (s.w+1)%s.size == s.r {
		s.condFull.Wait()
		if done {
			return errTimeout
		}
	}
	s.buffer[s.w] = val
	s.w++
	if s.w >= s.size {
		s.w = 0
	}
	s.condEmpty.Signal()
	return nil
}

func (s *SliceQueue[T]) Dequeue() (T, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for s.w == s.r {
		s.condEmpty.Wait()
	}
	s.mutex.Lock()
	result := s.buffer[s.r]
	s.buffer[s.r] = s.zero
	s.r++
	if s.r >= s.size {
		s.r = 0
	}
	s.condFull.Signal()
	return result, nil
}

type PriorityQueue[T any] struct {
}

type DelayQueue[T any] struct {
	dequeueChan chan struct{}
	lock        sync.RWMutex
	cond        sync.Cond
	PriorityQueue[T]
}

type Element[T any] struct {
}

func (elm *Element[T]) Deadline() time.Time {
	panic("implement me")
}

func (s *PriorityQueue[T]) Peek(ctx context.Context) (*Element[T], error) {
	panic("implement me")
}

func (s *PriorityQueue[T]) Dequeue(ctx context.Context) (*Element[T], error) {
	panic("implement me")
}

var (
	errEmptyQuery = errors.New("empty query")
	//errTimeout    = errors.New("timeout")
)

func NewDelayQueue[T any]() *DelayQueue[T] {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &DelayQueue[T]{
		dequeueChan: ch,
	}
}

func (s *DelayQueue[T]) Dequeue(ctx context.Context) (*Element[T], error) {
	select {
	case <-s.dequeueChan:
	case <-ctx.Done():
		return nil, errTimeout
	}
	defer func() {
		s.dequeueChan <- struct{}{}
	}()

	s.lock.Lock()
	defer s.lock.Unlock()
	select {
	case <-ctx.Done():
		return nil, errTimeout
	default:
	}

	for {
		elm, err := s.Peek(ctx)
		if err != errEmptyQuery && err != nil {
			return nil, err
		}
		var delay time.Duration
		if err == nil {
			now := time.Now()
			if elm.Deadline().Before(now) {
				elm, err = s.PriorityQueue.Dequeue(ctx)
				if err != nil {
					return nil, err
				}
				return elm, nil
			}
			delay = now.Sub(elm.Deadline())
		}
		go func() {
			if delay > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, delay)
				defer cancel()
			}
			<-ctx.Done()
			s.cond.Signal()
		}()
		s.cond.Wait()
	}
}

func (s *DelayQueue[T]) Enqueue(ctx context.Context, val T) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	select {
	case <-ctx.Done():
		return errTimeout
	default:
	}

	err := s.PriorityQueue.Enqueue(ctx, val)
	if err != nil {
		return err
	}
	s.cond.Signal()
	return nil
}

func (s *PriorityQueue[T]) Enqueue(ctx context.Context, val T) error {
	panic("implement me")
}

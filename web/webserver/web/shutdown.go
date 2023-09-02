package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type Hook func(ctx context.Context) error

type Shutdown struct {
	closing      atomic.Bool
	reqCount     atomic.Uint64
	chanShutdown chan struct{}
}

func NewShutdown() *Shutdown {
	return &Shutdown{
		chanShutdown: make(chan struct{}, 1),
	}
}

func (s *Shutdown) ShutdownFilterBuilder(next Filter) Filter {
	return func(c *Context) {
		if s.closing.Load() {
			c.W.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		s.reqCount.Add(1)
		next(c)
		n := s.reqCount.Add(^uint64(0))
		if s.closing.Load() && n == 0 {
			s.chanShutdown <- struct{}{}
		}
	}
}

func (s *Shutdown) TerminateAndWaitForShutdown(ctx context.Context) error {
	s.closing.Store(true)
	if s.reqCount.Load() == 0 {
		return nil
	}
	select {
	case <-ctx.Done():
		fmt.Println("time out terminating server")
		return errors.New("time out terminating server")
	case <-s.chanShutdown:
		fmt.Println("server terminated")
	}
	return nil
}

func (s *Shutdown) WaitForShutdown(hooks ...Hook) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, ShutdownSignals...)
	select {
	case sig := <-signals:
		fmt.Printf("get signal %v, application will shutdown\n", sig)
		time.AfterFunc(time.Minute, func() {
			fmt.Println("shutdown gracefully timeout, application will shutdown immediately")
			os.Exit(1)
		})
		for _, hook := range hooks {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			if err := hook(ctx); err != nil {
				fmt.Printf("failed to hook: %v\n", err)
			}
			cancel()
		}
		os.Exit(0)
	}
}

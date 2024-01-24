package errgroup

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestErrGroup_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	eg, egCtx := errgroup.WithContext(ctx)
	const count = 8
	for i := 0; i < count; i++ {
		eg.Go(func() error {
			return nil
		})
	}
	<-egCtx.Done()
	require.ErrorIs(t, context.DeadlineExceeded, egCtx.Err())
}

func TestErrGroup_Failed(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancel()
	eg, egCtx := errgroup.WithContext(ctx)
	const count = 8
	for i := 0; i < count; i++ {
		index := i
		eg.Go(func() error {
			if index == count/2 {
				return errors.New("failed")
			}
			return nil
		})
	}
	<-egCtx.Done()
	require.ErrorIs(t, context.Canceled, egCtx.Err())
}

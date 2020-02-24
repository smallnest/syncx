package syncx

import (
	"context"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

// Any is used to execute batch tasks anf wait some tasks finish.
// It can't be reused.
type Any struct {
	n     int
	least uint64
	p     uint32
	ctx   context.Context
	sema  *semaphore.Weighted

	finished uint64
	closedCh chan struct{}

	errs     []error
	errCount uint64
}

// NewAny creates a batch to execute n tasks and wait least tasks finish.
func NewAny(ctx context.Context, n int, least int) *Any {
	return &Any{
		ctx:      ctx,
		n:        n,
		least:    uint64(least),
		closedCh: make(chan struct{}),
		errs:     make([]error, n),
	}
}

// NewAnyWithParallel creates a batch to execute n tasks in max p goroutine and wait least tasks finish.
func NewAnyWithParallel(ctx context.Context, n int, p uint32, least int) *Any {
	return &Any{
		ctx:      ctx,
		n:        n,
		least:    uint64(least),
		closedCh: make(chan struct{}),
		p:        p,
		sema:     semaphore.NewWeighted(int64(p)),
		errs:     make([]error, n),
	}
}

// Go executes the ith task.
// Only call this method for one index one time.
func (b *Any) Go(i int, task func(ctx context.Context) error) {
	if b.sema != nil {
		err := b.sema.Acquire(b.ctx, 1)
		if err != nil {
			b.errs[i] = err
			atomic.AddUint64(&b.errCount, 1)
			b.finish()
			return
		}
	}
	go func() {
		defer b.finish()

		if b.sema != nil {
			defer b.sema.Release(1)
		}

		if err := task(b.ctx); err != nil {
			b.errs[i] = err
			atomic.AddUint64(&b.errCount, 1)
		}
	}()
}

func (b *Any) finish() {
	n := atomic.AddUint64(&b.finished, 1)
	if n == b.least {
		close(b.closedCh)
	}
}

// Wait blocks until least tasks have been done or canceled.
// if some tasks return errors, this method returns the count of errors and result of each task.
//
// Must be called after call Go methods.
func (b *Any) Wait() (errCount uint64, errs []error) {
	<-b.closedCh

	return b.errCount, b.errs
}

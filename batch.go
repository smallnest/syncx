package syncx

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

// Batch is used to execute batch tasks.
// It can't be reused.
type Batch struct {
	n    uint64
	gon  uint64
	p    uint32
	ctx  context.Context
	wg   sync.WaitGroup
	sema *semaphore.Weighted

	errs     []error
	errCount uint64
}

// NewBatch creates a batch to execute n tasks.
func NewBatch(ctx context.Context, n uint64) *Batch {
	return &Batch{
		ctx:  ctx,
		n:    n,
		errs: make([]error, n),
	}
}

// NewBatchWithParallel creates a batch to execute n tasks in max p goroutine.
func NewBatchWithParallel(ctx context.Context, n uint64, p uint32) *Batch {
	return &Batch{
		ctx:  ctx,
		n:    n,
		p:    p,
		sema: semaphore.NewWeighted(int64(p)),
		errs: make([]error, n),
	}
}

// Go executes the ith task.
// Only call this method for one index one time.
func (b *Batch) Go(i int, task func(ctx context.Context) error) {
	b.wg.Add(1)
	atomic.AddUint64(&b.gon, 1)

	if b.sema != nil {
		err := b.sema.Acquire(b.ctx, 1)
		if err != nil {
			b.errs[i] = err
			atomic.AddUint64(&b.errCount, 1)
			return
		}
	}
	go func() {
		defer b.wg.Done()
		if b.sema != nil {
			defer b.sema.Release(1)
		}

		if err := task(b.ctx); err != nil {
			b.errs[i] = err
			atomic.AddUint64(&b.errCount, 1)
		}
	}()
}

// Wait blocks until all tasks have been done or canceled.
// if some tasks return errors, this method returns the count of errors and result of each task.
// If some tasks have not added into the executer, this method returns the error.
//
// Must be called after call Go methods.
func (b *Batch) Wait() (errCount uint64, errs []error, err error) {
	n := atomic.LoadUint64(&b.gon)
	if n != b.n {
		return 0, nil, fmt.Errorf("some tasks have not been added. expected %d but current %d", b.n, n)
	}
	b.wg.Wait()

	return b.errCount, b.errs, nil
}

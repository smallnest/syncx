package syncx

import (
	"context"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

func TestAny(t *testing.T) {
	ctx := context.Background()
	any := NewAny(ctx, 100, 5)

	var count = uint64(0)

	for i := 0; i < 100; i++ {
		any.Go(i, func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				atomic.AddUint64(&count, 100)
				time.Sleep(time.Duration(rand.Int63n(1000)) * time.Millisecond)
				return nil
			}
		})
	}

	errCount, errs := any.Wait()
	if errCount != 0 {
		t.Fatalf("some goroutines return errors %d: %v", errCount, errs)
	}

	c := atomic.LoadUint64(&count)
	if c < 500 {
		t.Fatalf("expected >=500 but got %d", c)
	}
}

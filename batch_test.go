package syncx

import (
	"context"
	"sync/atomic"
	"testing"
)

func TestBatch(t *testing.T) {
	ctx := context.Background()
	batch := NewBatch(ctx, 3)

	var count = uint64(0)

	for i := 0; i < 3; i++ {
		batch.Go(i, func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				atomic.AddUint64(&count, 100)
				return nil
			}
		})
	}

	errCount, errs := batch.Wait()
	if errCount != 0 {
		t.Fatalf("some goroutines return errors %d: %v", errCount, errs)
	}

	if count != 300 {
		t.Fatalf("expected 300 but got %d", count)
	}
}

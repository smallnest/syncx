package syncx

import (
	"context"
	"strconv"
	"sync"
	"testing"
)

func TestToken(t *testing.T) {
	token := NewToken(3)

	// test 0 --> 1 ——> 2 --> 0 --> 1 --> 2 --> 0

	var wg sync.WaitGroup // 7 steps
	wg.Add(7)

	ctx, cancel := context.WithCancel(context.Background())

	var result string

	for i := 0; i < 3; i++ {
		i := i
		id := strconv.Itoa(i)
		go func() {
			for {
				err := token.Accquire(ctx, i)
				if err != nil {
					return
				}
				result += id
				err = token.Handoff(ctx, (i+1)%3) // pass to the next
				wg.Done()
				if err != nil {
					return
				}
			}
		}()
	}

	// begin
	err := token.Handoff(ctx, 0)
	if err != nil {
		t.Fatalf("failed to start the first step: %v", err)
	}

	wg.Wait()
	cancel()

	if result != "0120120" {
		t.Fatalf("expect steps 0120120 but got %s", result)
	}

}

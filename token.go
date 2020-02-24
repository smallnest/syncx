package syncx

import (
	"context"
	"math/rand"
)

// Token is a token shared by n participants, and only one participant can do stuff at the same time.
// Participants must accquire the otoen before they do stuff. They will be blocked if the token by another participant.
// After a participant has finished its work, it can handoff the token to an determined participant.
type Token struct {
	n     int
	chans []chan struct{}
}

// NewToken creates a new Token with n participants.
func NewToken(n int) *Token {
	chans := make([]chan struct{}, n)
	for i := 0; i < n; i++ {
		chans[i] = make(chan struct{}, 1)
	}

	return &Token{
		n:     n,
		chans: chans,
	}
}

// Accquire accquires the token by current participant(selfIdx).
func (t *Token) Accquire(ctx context.Context, selfIdx int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.chans[selfIdx]:
		return nil
	}
}

// Handoff passes the token to another participant(otherIdx).
func (t *Token) Handoff(ctx context.Context, otherIdx int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.chans[otherIdx] <- struct{}{}:
		return nil
	}
}

// Rand gives the token to a participant randomly.
func (t *Token) Rand(ctx context.Context) error {
	i := rand.Intn(t.n)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case t.chans[i] <- struct{}{}:
		return nil
	}
}

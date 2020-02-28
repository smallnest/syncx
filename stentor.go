package syncx

import (
	"sync"
)

// Stentor implments pub-sub pattern.
type Stentor struct {
	mu              sync.RWMutex
	bufSize         int
	subscribers     map[<-chan G]chan G
	FailureCallback func(<-chan G, G)
}

// NewStentor creates a Stentor.
func NewStentor(bufSize int) *Stentor {
	return &Stentor{
		bufSize:     bufSize,
		subscribers: make(map[<-chan G]chan G),
	}
}

// Subscribe subscribes this Stentor.
func (s *Stentor) Subscribe() <-chan G {
	ch := make(chan G, s.bufSize)
	s.mu.Lock()
	s.subscribers[ch] = ch
	s.mu.Unlock()

	return ch
}

// Unsubscribe unsubscribes this Stentor.
func (s *Stentor) Unsubscribe(sr <-chan G) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.unsubscribe(sr)
}

func (s *Stentor) unsubscribe(sr <-chan G) bool {
	ch, existed := s.subscribers[sr]
	if existed {
		delete(s.subscribers, sr)
		close(ch)
	}
	return existed
}

// Broadcast broadcasts an event.
func (s *Stentor) Broadcast(g G) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sr := range s.subscribers {
		select {
		case sr <- g:
		default:
			if s.FailureCallback != nil {
				s.FailureCallback(sr, g)
			}
		}
	}
}

// Count returns count of subscribers.
func (s *Stentor) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}

// Reset clears all subscribers.
func (s *Stentor) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for sr := range s.subscribers {
		s.unsubscribe(sr)
	}
	s.subscribers = make(map[<-chan G]chan G)
}

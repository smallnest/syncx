package syncx

import "testing"

func TestStentor(t *testing.T) {
	s := NewStentor(100)

	ob1 := s.Subscribe()
	ob2 := s.Subscribe()

	for i := 0; i < 10; i++ {
		s.Broadcast(i)
	}

	if len(ob1) != 10 {
		t.Fatalf("expect 10 but got %d", len(ob1))
	}
	if len(ob2) != 10 {
		t.Fatalf("expect 10 but got %d", len(ob2))
	}

	for i := 0; i < 1000; i++ {
		s.Broadcast(i)
	}

	if len(ob1) != 100 {
		t.Fatalf("expect 100 but got %d", len(ob1))
	}
	if len(ob2) != 100 {
		t.Fatalf("expect 100 but got %d", len(ob2))
	}

	count := s.Count()
	if count != 2 {
		t.Fatalf("expect 2 but got %d", count)
	}

	s.Reset()
	count = s.Count()
	if count != 0 {
		t.Fatalf("expect 0 but got %d", count)
	}
}

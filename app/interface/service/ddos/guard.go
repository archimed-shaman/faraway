package ddos

import (
	"container/ring"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type DDoSGuard struct {
	window *ring.Ring

	current  atomic.Int64
	windowed atomic.Int64

	initOnce sync.Once
}

func NewGuard(window time.Duration) *DDoSGuard {
	n := int(window.Seconds()) - 1
	if n < 0 {
		n = 0
	}

	return &DDoSGuard{
		window:   ring.New(n),
		current:  atomic.Int64{},
		windowed: atomic.Int64{},
		initOnce: sync.Once{},
	}
}

func (s *DDoSGuard) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.reset()
		}
	}
}

func (s *DDoSGuard) IncRate(_ context.Context) (int64, error) {
	return s.windowed.Load() + s.current.Add(1), nil
}

func (s *DDoSGuard) reset() {
	lastSec := s.current.Swap(0)

	if s.window.Len() > 0 {
		if value, ok := s.window.Value.(int64); ok {
			s.windowed.Add(-value) // remove the current sum, it will be replaced
		}

		s.window.Value = lastSec
		s.windowed.Add(lastSec)
		s.window = s.window.Next()
	}
}

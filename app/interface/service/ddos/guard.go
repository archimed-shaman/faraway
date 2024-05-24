package ddos

import (
	"context"
	"sync/atomic"
)

type DDoSGuard struct {
	count atomic.Int64
}

func NewGuard() *DDoSGuard {
	return &DDoSGuard{
		count: atomic.Int64{},
	}
}

func (s *DDoSGuard) IncRate(ctx context.Context) (int64, error) {
	return s.count.Add(1), nil
}

func (s *DDoSGuard) DecRate(ctx context.Context) (int64, error) {
	return s.count.Add(-1), nil
}

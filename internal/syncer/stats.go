package syncer

import (
	"sync/atomic"
)

// Stats synchronizes block cursors values across the head and janitor.
type Stats struct {
	headCursor atomic.Uint64
	lowerBound atomic.Uint64
}

func (s *Stats) UpdateHeadCursor(val uint64) {
	s.headCursor.Store(val)
}

func (s *Stats) GetHeadCursor() uint64 {
	return s.headCursor.Load()
}

func (s *Stats) UpdateLowerBound(val uint64) {
	s.lowerBound.Store(val)
}

func (s *Stats) GetLowerBound() uint64 {
	return s.lowerBound.Load()
}

package semafore

import "log"

type Semaphore struct {
	ch  chan struct{}
	max int
}

func NewSemaphore(max int) *Semaphore {
	return &Semaphore{
		ch:  make(chan struct{}, max),
		max: max,
	}
}

func (s *Semaphore) Acquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		log.Printf("429: too many concurrent requests (%d/%d)", len(s.ch), s.max)
		return false
	}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func (s *Semaphore) Usage() int {
	return len(s.ch)
}

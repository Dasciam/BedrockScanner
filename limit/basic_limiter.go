package limit

import (
	"sync"
	"time"
)

type BasicLimiter struct {
	last  time.Time
	delay time.Duration

	mu sync.Mutex
}

func NewBasicLimiter(perSecond int) *BasicLimiter {
	return &BasicLimiter{
		last:  time.Now(),
		delay: time.Second / time.Duration(perSecond),
	}
}

func (b *BasicLimiter) Increment() {
	b.mu.Lock()
	defer b.mu.Unlock()

	time.Sleep(time.Until(b.last))
	b.last = b.last.Add(b.delay)
}

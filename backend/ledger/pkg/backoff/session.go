package backoff

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Session struct {
	heapCtrl       *heapController
	key            string
	cfg            *sessionOptions
	nextAllowedAt  time.Time
	mu             sync.RWMutex
	expiredAt      time.Time
	expired        int32
	counter        int
	backoffCounter int
}

func newSession(heapCtrl *heapController, key string, opts *sessionOptions) *Session {
	s := &Session{
		heapCtrl:      heapCtrl,
		key:           key,
		cfg:           opts,
		nextAllowedAt: time.Now(),
		expiredAt:     time.Now().Add(opts.Ttl),
	}

	heapCtrl.add(s)
	return s
}

func (s *Session) isExpired() bool {
	return time.Now().After(s.expiredAt)
}

func (s *Session) AddBackoff() bool {
	if atomic.LoadInt32(&s.expired) == 1 {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	mul := time.Duration(math.Pow(s.cfg.Multiplier, float64(s.backoffCounter)))

	s.backoffCounter++

	diffTime := time.Duration(s.cfg.InitialInterval * mul)
	if diffTime > s.cfg.MaxInterval {
		diffTime = s.cfg.MaxInterval
	}

	s.nextAllowedAt = time.Now().Add(diffTime)
	s.expiredAt = time.Now().Add(s.cfg.Ttl)
	s.heapCtrl.update(s)

	return true
}

func (s *Session) Counter() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.counter
}

func (s *Session) AddCounter(value ...int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val := 1
	if len(value) > 0 {
		val = value[0]
	}

	s.counter += val
	s.expiredAt = time.Now().Add(s.cfg.Ttl)
	s.heapCtrl.update(s)
}

func (s *Session) Key() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.key
}

func (s *Session) IsAllowed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().After(s.nextAllowedAt) || time.Now().Equal(s.nextAllowedAt)
}

func (s *Session) Reset() {
	if atomic.LoadInt32(&s.expired) == 1 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter = 0
	s.backoffCounter = 0
	s.nextAllowedAt = time.Now()
	s.expiredAt = time.Now().Add(s.cfg.Ttl)
	s.heapCtrl.update(s)
}

func (s *Session) NextAllowed() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nextAllowedAt
}

func (s *Session) NextAllowedUntilSeconds() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	diff := time.Until(s.nextAllowedAt).Seconds()
	if diff < 0 {
		diff = 0
	}

	return int64(math.Ceil(diff))
}

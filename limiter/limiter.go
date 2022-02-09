package limiter

import (
	"sync"
	"time"
)

func NewLimiter(interval time.Duration, maxCount int) *Limiter {
	return &Limiter{
		interval:  interval,
		maxCount:  maxCount,
		startTime: time.Now(),
	}
}

type Limiter struct {
	interval  time.Duration
	maxCount  int
	lock      sync.Mutex
	startTime time.Time
	curCount  int
}

func (l *Limiter) Check() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	nowTime := time.Now()
	if nowTime.After(l.startTime.Add(l.interval)) {
		l.startTime = nowTime
		l.curCount = 0
	}
	l.curCount++
	return l.curCount < l.maxCount
}

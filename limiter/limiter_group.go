package limiter

import (
	"sync"
	"time"
)

func NewLimiterGroup() *LimiterGroup {
	return &LimiterGroup{
		limiterMap: make(map[interface{}]*Limiter),
	}
}

type LimiterGroup struct {
	limiterMap map[interface{}]*Limiter
	lock       sync.RWMutex
}

func (lm *LimiterGroup) Add(key interface{}, interval time.Duration, maxCount int) {
	l := NewLimiter(interval, maxCount)

	lm.lock.Lock()
	defer lm.lock.Unlock()

	lm.limiterMap[key] = l
}

func (lm *LimiterGroup) Check(key interface{}) bool {
	l := lm.get(key)

	return l.Check()
}

func (lm *LimiterGroup) get(key interface{}) *Limiter {
	lm.lock.RLock()
	defer lm.lock.RUnlock()

	return lm.limiterMap[key]
}

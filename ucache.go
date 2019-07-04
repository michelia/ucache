package ucache

import (
	"sync"
	"time"
)

type setExpired struct {
	expiration time.Duration
	interval   time.Duration
	items      map[string]int64
	mu         sync.RWMutex
}

/*
NewSetExpired 创建带有过期时间的集合(类似与Python Set)

expiration: key的过期时间.
interval: 清理过期key的间隔时间.
*/
func NewSetExpired(expiration, interval time.Duration) *setExpired {
	s := &setExpired{
		expiration: expiration,
		interval:   interval,
		items:      make(map[string]int64),
	}
	go s.run()
	return s
}

// Has 是否包含当前key
func (s *setExpired) Has(k string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.items[k]; ok {
		return true
	}
	return false
}

// Add 添加当前key
func (s *setExpired) Add(k string) {
	s.mu.Lock()
	s.items[k] = time.Now().Add(s.expiration).Unix()
	s.mu.Unlock()
}

// run 定期清理过期key
func (s *setExpired) run() {
	ticker := time.NewTicker(s.interval)
	for {
		select {
		case <-ticker.C:
			s.deleteExpired()
		}
	}
}

// deleteExpired 清理过期key
func (s *setExpired) deleteExpired() {
	s.mu.Lock()
	now := time.Now().Unix()
	for k, v := range s.items {
		if now > v {
			delete(s.items, k)
		}
	}
	s.mu.Unlock()
}

package main

import (
	"fmt"
	"sync"
)

// Solution 2: Thread-Safe Map using sync.RWMutex + regular map
// RWMutex allows multiple concurrent readers but only one writer at a time.
// We use Lock() for writes and RLock() for reads to protect the map.

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[string]int)}
}

func (s *SafeMap) Set(key string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap) Get(key string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.m[key]
	return val, ok
}

func main() {
	safeMap := NewSafeMap()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			safeMap.Set("key", key)
		}(i)
	}

	wg.Wait()

	value, _ := safeMap.Get("key")
	fmt.Printf("Value: %d\n", value)
}

package main

import (
	"fmt"
	"sync"
)

// Solution 1: Thread-Safe Map using sync.Map
// sync.Map is a built-in concurrent map that handles locking internally.
// No external mutex is needed — safe for concurrent reads and writes.

func main() {
	var sm sync.Map
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			sm.Store("key", key)
		}(i)
	}

	wg.Wait()

	value, _ := sm.Load("key")
	fmt.Printf("Value: %v\n", value)
}

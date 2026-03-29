package main

import (
	"fmt"
	"sync"
)

// Solution 1: Fix Concurrent Counter using sync.Mutex
//
// Why the original counter is not 1000:
// The counter++ operation is NOT atomic — it consists of three steps:
// READ the value, INCREMENT it, and WRITE it back. When multiple goroutines
// execute these steps simultaneously without synchronization, some increments
// are lost because goroutines overwrite each other's results (this is a DATA RACE).
//
// Run with: go run -race problem2_mutex.go

func main() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(counter) // Always prints 1000
}

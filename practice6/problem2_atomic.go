package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Solution 2: Fix Concurrent Counter using sync/atomic
//
// atomic.AddInt64 performs the increment as a single indivisible CPU instruction,
// so no goroutine can interrupt it mid-operation — no mutex needed.
//
// Run with: go run -race problem2_atomic.go

func main() {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println(atomic.LoadInt64(&counter)) // Always prints 1000
}

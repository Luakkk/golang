package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Stage 1: startServer simulates a server sending metrics in real-time.
// It sends a new metric at a random interval (0-500ms) until the context is done.
func startServer(ctx context.Context, name string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(rand.Intn(500)) * time.Millisecond):
				out <- fmt.Sprintf("[%s] metric: %d", name, rand.Intn(100))
			}
		}
	}()
	return out
}

// FanIn merges multiple input channels into a single output channel.
// It uses sync.WaitGroup to close the result channel only after all
// input channels are fully drained and closed.
// Context support ensures the fan-in also stops when ctx is cancelled.
func FanIn(ctx context.Context, channels ...<-chan string) <-chan string {
	merged := make(chan string)
	var wg sync.WaitGroup

	// For each input channel, start a goroutine that forwards its values
	// to the merged channel.
	forward := func(ch <-chan string) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}
				select {
				case merged <- val:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go forward(ch)
	}

	// Close the merged channel once all forwarder goroutines are done.
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch1 := startServer(ctx, "Alpha")
	ch2 := startServer(ctx, "Beta")
	ch3 := startServer(ctx, "Gamma")

	ch4 := FanIn(ctx, ch1, ch2, ch3)

	for val := range ch4 {
		fmt.Println(val)
	}
}

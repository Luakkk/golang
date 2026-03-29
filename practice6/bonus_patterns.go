package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================
// BONUS: Concurrency Patterns
// ============================================================
// This file demonstrates two classic Go concurrency patterns:
//   1. Worker Pool  — distribute tasks across N workers
//   2. Pipeline     — chain processing stages via channels
// Run with: go run -race bonus_patterns.go
// ============================================================

// ------------------------------------------------------------
// Pattern 1: Worker Pool
// ------------------------------------------------------------
// A fixed set of goroutines (workers) reads jobs from a shared
// channel and sends results to a results channel.
// This bounds memory usage and parallelism for CPU-heavy tasks.

type Job struct {
	ID    int
	Value int
}

type Result struct {
	Job    Job
	Output int
}

// worker processes jobs from the jobs channel and sends results
// to the results channel. It stops when jobs is closed.
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Simulate some work (squaring the value)
		time.Sleep(10 * time.Millisecond)
		results <- Result{Job: job, Output: job.Value * job.Value}
		fmt.Printf("[Worker %d] processed job %d: %d^2 = %d\n", id, job.ID, job.Value, job.Value*job.Value)
	}
}

func runWorkerPool() {
	fmt.Println("\n=== Worker Pool Pattern ===")

	const numWorkers = 3
	const numJobs = 9

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)

	var wg sync.WaitGroup

	// Start workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- Job{ID: j, Value: j}
	}
	close(jobs) // signal workers that no more jobs are coming

	// Close results once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	total := 0
	for res := range results {
		total += res.Output
	}
	fmt.Printf("Worker Pool done. Sum of squares 1..%d = %d\n", numJobs, total)
}

// ------------------------------------------------------------
// Pattern 2: Pipeline
// ------------------------------------------------------------
// Data flows through a series of stages, each implemented as a
// goroutine that reads from an input channel and writes to an
// output channel.  Stages are composed by connecting channels.

// Stage 1: generate integers 1..n
func generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Stage 2: square each integer
func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Stage 3: filter — only pass through even numbers
func filterEven(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

func runPipeline() {
	fmt.Println("\n=== Pipeline Pattern ===")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Build pipeline: generate → square → filterEven
	nums := generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squares := square(ctx, nums)
	evens := filterEven(ctx, squares)

	fmt.Print("Even squares of 1..10: ")
	for val := range evens {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}

func main() {
	runWorkerPool()
	runPipeline()
}

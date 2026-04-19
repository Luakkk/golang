package main

// ADVANCED bonus: RetryExchangeService wraps ExchangeService and retries
// on any error up to maxRetry times with an optional delay between attempts.

import (
	"fmt"
	"time"
)

// RetryExchangeService wraps an ExchangeService and retries failed GetRate calls.
type RetryExchangeService struct {
	svc      *ExchangeService
	maxRetry int
	delay    time.Duration
}

// NewRetryExchangeService creates a RetryExchangeService.
// maxRetry is the total number of attempts (e.g. 3 means: try 1, retry 2, retry 3).
// delay is the pause between attempts (use 0 in tests for speed).
func NewRetryExchangeService(svc *ExchangeService, maxRetry int, delay time.Duration) *RetryExchangeService {
	return &RetryExchangeService{svc: svc, maxRetry: maxRetry, delay: delay}
}

// GetRate calls the underlying GetRate up to maxRetry times.
// Returns the first successful result, or wraps the last error if all retries fail.
func (r *RetryExchangeService) GetRate(from, to string) (float64, error) {
	var err error
	for attempt := 1; attempt <= r.maxRetry; attempt++ {
		var rate float64
		rate, err = r.svc.GetRate(from, to)
		if err == nil {
			return rate, nil
		}
		if attempt < r.maxRetry && r.delay > 0 {
			time.Sleep(r.delay)
		}
	}
	return 0, fmt.Errorf("all %d retries failed: %w", r.maxRetry, err)
}

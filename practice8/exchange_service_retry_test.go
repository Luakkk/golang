package main

// ADVANCED bonus: Tests for RetryExchangeService.
// Scenarios:
//   1. First request fails, second succeeds.
//   2. All retries fail — error wraps "all N retries failed".

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRetryGetRate_FirstFailSecondSuccess: first request causes a network error
// (server panics → connection closed), second request returns a valid response.
func TestRetryGetRate_FirstFailSecondSuccess(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			// Simulate server failure on first attempt: panic closes the connection,
			// and the client receives a network/EOF error — exactly what retry handles.
			panic("simulated failure on first attempt")
		}
		// Second attempt: respond normally.
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RateResponse{Base: "USD", Target: "EUR", Rate: 0.92})
	}))
	defer server.Close()

	inner := &ExchangeService{
		BaseURL: server.URL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
	svc := NewRetryExchangeService(inner, 3, 0) // 3 attempts, no delay in tests

	rate, err := svc.GetRate("USD", "EUR")

	require.NoError(t, err)
	assert.Equal(t, 0.92, rate)
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount),
		"expected exactly 2 server calls: 1 failed + 1 success")
}

// TestRetryGetRate_AllRetriesFail: every attempt fails; the returned error
// must mention "all 3 retries failed".
func TestRetryGetRate_AllRetriesFail(t *testing.T) {
	var callCount int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		panic("always fails")
	}))
	defer server.Close()

	inner := &ExchangeService{
		BaseURL: server.URL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
	svc := NewRetryExchangeService(inner, 3, 0)

	rate, err := svc.GetRate("USD", "EUR")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "all 3 retries failed")
	assert.Equal(t, float64(0), rate)
	assert.Equal(t, int32(3), atomic.LoadInt32(&callCount),
		"expected exactly 3 server calls (all failed)")
}

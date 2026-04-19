package main

// Task 3: Tests for GetRate using net/http/httptest to mock the external API.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper: build an ExchangeService pointing at the given test server URL.
func newTestExchangeService(serverURL string) *ExchangeService {
	return &ExchangeService{
		BaseURL: serverURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// 1. Successful scenario
func TestGetRate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(RateResponse{Base: "USD", Target: "EUR", Rate: 0.92})
	}))
	defer server.Close()

	svc := newTestExchangeService(server.URL)
	rate, err := svc.GetRate("USD", "EUR")

	require.NoError(t, err)
	assert.Equal(t, 0.92, rate)
}

// 2. API Business Error: server returns 400 / 404 with {"error": "invalid currency pair"}
func TestGetRate_APIBusinessError_400(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RateResponse{ErrorMsg: "invalid currency pair"})
	}))
	defer server.Close()

	svc := newTestExchangeService(server.URL)
	rate, err := svc.GetRate("INVALID", "PAIR")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid currency pair")
	assert.Equal(t, float64(0), rate)
}

func TestGetRate_APIBusinessError_404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(RateResponse{ErrorMsg: "invalid currency pair"})
	}))
	defer server.Close()

	svc := newTestExchangeService(server.URL)
	rate, err := svc.GetRate("INVALID", "PAIR")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid currency pair")
	assert.Equal(t, float64(0), rate)
}

// 3. Malformed JSON: server returns non-JSON text
func TestGetRate_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("this is NOT valid JSON {{{"))
	}))
	defer server.Close()

	svc := newTestExchangeService(server.URL)
	rate, err := svc.GetRate("USD", "EUR")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode error")
	assert.Equal(t, float64(0), rate)
}

// 4. Slow Response / Timeout: server sleeps longer than client timeout
func TestGetRate_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // much longer than our 100 ms timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	svc := &ExchangeService{
		BaseURL: server.URL,
		Client:  &http.Client{Timeout: 100 * time.Millisecond},
	}
	rate, err := svc.GetRate("USD", "EUR")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
	assert.Equal(t, float64(0), rate)
}

// 5. Server Panic → net/http recovers the panic and closes the connection,
//    the client gets a network/EOF error.
func TestGetRate_ServerPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("simulated server panic")
	}))
	defer server.Close()

	svc := newTestExchangeService(server.URL)
	rate, err := svc.GetRate("USD", "EUR")

	require.Error(t, err)
	assert.Equal(t, float64(0), rate)
}

// 6. Empty Body: server returns 200 with no body → JSON decode fails
func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// intentionally write nothing
	}))
	defer server.Close()

	svc := newTestExchangeService(server.URL)
	rate, err := svc.GetRate("USD", "EUR")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode error")
	assert.Equal(t, float64(0), rate)
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sumitwarrior7/social/internal/auth"
	"github.com/Sumitwarrior7/social/internal/ratelimiter"
	"github.com/Sumitwarrior7/social/internal/store"
	"github.com/Sumitwarrior7/social/internal/store/cache"
	"go.uber.org/zap"
)

// func newTestApplication(t *testing.T) *application {
// 	t.Helper()

// 	logger := zap.NewNop().Sugar()
// 	mockStore := store.NewMockStore()
// 	mockCacheStore := cache.NewMockStore()
// 	testAuth := &auth.TestAuthenticator{}

// 	return &application{
// 		logger:        logger,
// 		store:         mockStore,
// 		cacheStorage:  mockCacheStore,
// 		authenticator: testAuth,
// 	}
// }

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// Uncomment to enable logs
	// logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected error response %d. Got %d", expected, actual)
	}
}

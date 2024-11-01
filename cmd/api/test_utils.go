package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/damarteplok/social/internal/auth"
	"github.com/damarteplok/social/internal/store"
	"github.com/damarteplok/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	os.Setenv("ENV", "test")
	defer func() {
		os.Unsetenv("ENV")
	}()

	// logger := zap.NewNop().Sugar()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}

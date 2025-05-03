package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"a1337b04rd/internal/adapter/http/middleware"
)

func TestLoggerMiddleware(t *testing.T) {
	// флаг для проверки, вызвался ли хендлер
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rr := httptest.NewRecorder()

	mw := middleware.LoggerMiddleware(handler)
	mw.ServeHTTP(rr, req)

	if !called {
		t.Errorf("handler was not called")
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v want %v", status, http.StatusOK)
	}
}

package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {

	output := strings.Builder{}
	logger := log.New(&output, "", 0)

	mw := NewLogger(logger)
	timeout := 10 * time.Millisecond

	handler := mw(
		func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(timeout)
		},
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	handler(w, r)

	logRecord := output.String()
	if !strings.Contains(logRecord, r.RemoteAddr) {
		t.Errorf("expected log record '%s' to contain remote addres %s", logRecord, r.RemoteAddr)
	}

	if !strings.Contains(logRecord, r.URL.Path) {
		t.Errorf("expected log record '%s' to contain path %s", logRecord, r.URL.Path)
	}

	if !strings.Contains(logRecord, fmt.Sprintf("%d", w.Code)) {
		t.Errorf("expected log record '%s' to contain response code %d", logRecord, w.Code)
	}

	if !strings.Contains(logRecord, timeout.String()) {
		t.Errorf("expected log record '%s' to contain execution time %s", logRecord, timeout.String())
	}

}

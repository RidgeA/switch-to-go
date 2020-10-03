package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	Index(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got: %d", http.StatusOK, w.Code)
	}
}

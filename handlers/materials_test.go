package handlers

import (
	"context"
	"errors"
	"github.com/RidgeA/switch-to-go/db/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockMaterialsGetter struct {
	response    []repository.Material
	responseErr error
}

func (m mockMaterialsGetter) GetMaterials(ctx context.Context) ([]repository.Material, error) {
	if m.response != nil {
		return m.response, nil
	}
	return nil, m.responseErr
}

func TestMaterials_Success(t *testing.T) {

	r := httptest.NewRequest("GET", "/materials", nil)
	w := httptest.NewRecorder()
	url := "https://example.com"
	handler := Materials(mockMaterialsGetter{
		response: []repository.Material{{Url: url}},
	})

	handler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got: %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, url) {
		t.Errorf("expected body to contain url '%s'", url)
	}
}

func TestMaterials_Error(t *testing.T) {
	r := httptest.NewRequest("GET", "/materials", nil)
	w := httptest.NewRecorder()
	expectedBody := "can not get list of materials: some error"
	responseErr := errors.New("some error")
	handler := Materials(mockMaterialsGetter{
		responseErr: responseErr,
	})

	handler(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got: %d", http.StatusInternalServerError, w.Code)
	}

	if w.Body.String() != expectedBody {
		t.Errorf("expected body to be equal '%s'", expectedBody)
	}
}

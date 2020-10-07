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

func TestMaterials2(t *testing.T) {

	type args struct {
		w      *httptest.ResponseRecorder
		r      *http.Request
		getter MaterialsGetter
	}
	type expected struct {
		status int
		body   string
	}
	tests := []struct {
		name string
		args args
		expected
	}{
		{
			name: "success",

			args: args{
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "/materials", nil),
				mockMaterialsGetter{
					response: []repository.Material{{Url: "https://example.com"},
					},
				},
			},

			expected: expected{
				http.StatusOK,
				"https://example.com",
			},
		},
		{
			name: "error",

			args: args{
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "/materials", nil),
				mockMaterialsGetter{
					responseErr: errors.New("some error"),
				},
			},

			expected: expected{
				http.StatusInternalServerError,
				"some error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			handler := Materials(tt.args.getter)

			handler(tt.args.w, tt.args.r)

			if tt.args.w.Code != tt.expected.status {
				t.Errorf("expected status %d, got: %d", tt.expected.status, tt.args.w.Code)
			}

			body := tt.args.w.Body.String()
			if !strings.Contains(body, tt.expected.body) {
				t.Errorf("expected body to contain url '%s'", tt.expected.body)
			}
		})
	}
}

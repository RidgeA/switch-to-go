package middlewares

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type (
	testCache map[string]interface{}
)

func (t testCache) Set(s string, i interface{}) {
	t[s] = i
}

func (t testCache) Get(s string) (interface{}, bool) {
	v, ok := t[s]
	return v, ok
}

func (t testCache) Del(s string) {
	delete(t, s)
}

func TestNewCacher_NewRequest(t *testing.T) {

	path := "/"
	method := "GET"
	respBody := []byte("hello")

	cache := testCache(map[string]interface{}{})
	cacher := NewCacher(cache)

	handler := cacher(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(respBody)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	handler(w, r)

	v, _ := cache.Get(fmt.Sprint(method, path))

	switch rce := v.(type) {
	case ResponseCacheEntry:
		if bytes.Compare(rce.Body, respBody) != 0 {
			t.Errorf("expect to get %s, got %s", respBody, rce.Body)
		}
	default:
		t.Errorf("expect to get from cache ResponseCacheEntry, got %T", v)
	}

}

func TestNewCacher_CachedRequest(t *testing.T) {
	path := "/"
	method := "GET"
	rce := ResponseCacheEntry{
		Body:   []byte("hello"),
		Status: http.StatusCreated,
		Headers: map[string][]string{
			"Key": {"value"},
		},
	}
	cache := testCache(map[string]interface{}{
		"GET/": rce,
	})
	cacher := NewCacher(cache)

	handler := cacher(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("this should never be called"))
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	handler(w, r)

	if bytes.Compare(w.Body.Bytes(), rce.Body) != 0 {
		t.Errorf("expect to get in body %s, got %s", rce.Body, w.Body.String())
	}

	if w.Code != rce.Status {
		t.Errorf("expect to get status %d, got %d", rce.Status, w.Code)
	}

	if !reflect.DeepEqual(w.Header(), rce.Headers) {
		t.Errorf("expected to get cached headers")
	}
}

func TestNewCacher_CachedRequest_brokenCache(t *testing.T) {
	path := "/"
	method := "GET"
	expectedBody := []byte("this should be called due to broken cache record")
	cache := testCache(map[string]interface{}{
		"GET/": []byte("hello"),
	})
	cacher := NewCacher(cache)

	handler := cacher(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(expectedBody)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil)

	handler(w, r)

	if bytes.Compare(w.Body.Bytes(), expectedBody) != 0 {
		t.Errorf("expect to get %s, got %s", expectedBody, w.Body.String())
	}
}

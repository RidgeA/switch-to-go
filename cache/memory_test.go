package cache

import (
	"testing"
	"time"
)

func TestNewMemoryCache(t *testing.T) {
	d := time.Second
	got := NewMemoryCache(d)
	if got.storage == nil {
		t.Errorf("internal storage should be initialized")
	}
	if got.duration != d {
		t.Errorf("duration should be set")
	}
}

func TestMemoryCache_Get(t *testing.T) {
	d := time.Hour
	key := "key"
	value := "value"
	mc := NewMemoryCache(d)
	mc.Set(key, value)

	got, exists := mc.Get(key)

	if got != nil && !exists {
		t.Errorf("expected true if key exists")
	}

	if got == nil && !exists {
		t.Errorf("expected key to be in the cache")
	}

	gotStr := got.(string)
	if gotStr != value {
		t.Errorf("expected %s to be eqal %s", gotStr, value)
	}
}

func TestMemoryCache_Get_timeout(t *testing.T) {
	d := 100 * time.Millisecond
	key := "key"
	value := "value"
	mc := NewMemoryCache(d)
	mc.Set(key, value)

	got, _ := mc.Get(key)

	gotStr := got.(string)
	if gotStr != value {
		t.Errorf("expected %s to be eqal %s right after set", gotStr, value)
	}

	time.Sleep(110 * time.Millisecond)

	got2, exists := mc.Get(key)
	if got2 != nil || exists {
		t.Errorf("expect value not exits after specivied time, got: %s, %t", got2, exists)
	}
}

func TestMemoryCache_Del(t *testing.T) {
	d := 100 * time.Millisecond
	key := "key"
	value := "value"
	mc := NewMemoryCache(d)
	mc.Set(key, value)

	got, _ := mc.Get(key)

	gotStr := got.(string)
	if gotStr != value {
		t.Errorf("expected %s to be eqal %s right after set", gotStr, value)
	}

	mc.Del(key)

	got2, exists := mc.Get(key)
	if got2 != nil || exists {
		t.Errorf("expect value not exits after calling 'Del', got: %s, %t", got2, exists)
	}
}

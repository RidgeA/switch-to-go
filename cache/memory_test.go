package cache

import (
	"fmt"
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

func ExampleNewMemoryCache() {
	d := 100 * time.Millisecond
	key := "Hello"
	value := "World"

	cache := NewMemoryCache(d)

	// Save some key-value pair to cache
	cache.Set(key, value)

	// Get value by key
	getValue, exists := cache.Get(key)
	fmt.Printf("Exists: %t, value: %s\n", exists, getValue)

	// Delete value by key
	cache.Del(key)
	_, exists = cache.Get(key)
	fmt.Printf("Exists: %t\n", exists)

	// Output:
	// Exists: true, value: World
	// Exists: false
}

func BenchmarkMemoryCache_Get(b *testing.B) {
	d := time.Hour
	key := "Hello"
	value := "World"

	mc := NewMemoryCache(d)
	mc.Set(key, value)

	for i := 0; i < b.N; i++ {
		mc.Get(key)
	}
}
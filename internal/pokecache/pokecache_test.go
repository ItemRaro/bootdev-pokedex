package pokecache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestGetMissingKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	val, ok := cache.Get("https://example.com/missing")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
	if val != nil {
		t.Errorf("expected nil value, got %s", val)
		return
	}
}

func TestAddOverwrite(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("https://example.com", []byte("old"))
	cache.Add("https://example.com", []byte("new"))

	val, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}
	if string(val) != "new" {
		t.Errorf("expected new, got %s", val)
		return
	}
}

// TestReapLoopKeepsFreshEntries checks that entries survive until the interval passes.
func TestReapLoopKeepsFreshEntries(t *testing.T) {
	const interval = 100 * time.Millisecond
	cache := NewCache(interval)
	cache.Add("https://example.com", []byte("testdata"))

	time.Sleep(interval / 5)

	if _, ok := cache.Get("https://example.com"); !ok {
		t.Errorf("expected to find key before the interval")
		return
	}
}

// TestConcurrentAccess only proves something when run with "go test -race".
func TestConcurrentAccess(t *testing.T) {
	const goroutines = 100
	cache := NewCache(5 * time.Second)

	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Go(func() {
			key := fmt.Sprintf("https://example.com/%d", i)
			cache.Add(key, []byte("testdata"))
			cache.Get(key)
		})
	}
	wg.Wait()

	for i := range goroutines {
		key := fmt.Sprintf("https://example.com/%d", i)
		if _, ok := cache.Get(key); !ok {
			t.Errorf("expected to find key %s", key)
		}
	}
}

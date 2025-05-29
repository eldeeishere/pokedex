package pokecache

import (
	"testing"
	"time"
)

func TestCach(t *testing.T) {
	cases := []struct {
		key      string
		val      []byte
		expected bool
		timeout  time.Duration
	}{
		{
			key:      "hello",
			val:      []byte("world"),
			expected: true,
			timeout:  2 * time.Second,
		},
		{
			key:      "test2",
			val:      []byte("janepanemajim"),
			expected: false,
			timeout:  6 * time.Second,
		},
	}

	for _, c := range cases {
		cache := NewCache(4 * time.Second)
		cache.Add(c.key, c.val)
		time.Sleep(c.timeout)
		val, ok := cache.Get(c.key)
		if ok != c.expected {
			t.Errorf("Expected cache hit for key '%s', got %v", c.key, ok)
			continue
		}
		if c.expected {
			if string(val) != string(c.val) {
				t.Errorf("Expected value '%s' for key '%s', got '%s'", c.val, c.key, val)
			}
		}

	}

}

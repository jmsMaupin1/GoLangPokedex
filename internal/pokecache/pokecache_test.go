package pokecache

import (
	"testing"
	"time"
	"fmt"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second 
	cases := []struct{
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/test",
			val: []byte("heres another test"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)

			if !ok {
				t.Errorf("Expected to find key: %s", c.key)
				return
			}

			if string(c.val) != string(val) {
				t.Errorf("Expected to find value: %s", c.val)
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Second
	const waitTime = baseTime + 5 * time.Second

	cache := NewCache(baseTime)
	cache.Add("test.com", []byte("hello there bro"))

	_, ok := cache.Get("test.com")
	if !ok {
		t.Errorf("Expected to find key")
	}

	time.Sleep(waitTime)

	_, rok := cache.Get("test.com")
	if rok {
		t.Errorf("Expected key to be removed")
	}
}

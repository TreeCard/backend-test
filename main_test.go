package main

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	var count int

	key := "key"
	value := "value"
	getter := func() (string, time.Time, bool) {
		count++
		return value, time.Time{}, true
	}
	c := NewCache(getter, 10)

	check, ok := c.Get(key)
	if !ok {
		t.Error("getter failed")
	}
	if value != check {
		t.Errorf("key mismatched %s != %s", key, value)
	}
	if count != 1 {
		t.Errorf("count %d times", count)
	}

	// check get with a nil getter fails
	check, ok = c.Get(key)
	if !ok {
		t.Error("get shouldn't fail")
	}
	if check != value {
		t.Errorf("got %q expected %q", check, value)
	}
	if count != 1 {
		t.Errorf("count %d times", count)
	}
}

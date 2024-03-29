package cache_test

import (
	"github.com/treecard/backend-test/cache"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	var count int

	key := "key"
	value := "value"
	getter := func(_ string) (string, time.Time, bool) {
		count++
		return value, time.Time{}, true
	}
	c := cache.NewCache(getter, 10)

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

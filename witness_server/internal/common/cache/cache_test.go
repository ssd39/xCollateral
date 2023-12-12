package cache

import (
	"testing"
	"time"
)

func TestCache_GetAndSet(t *testing.T) {
	expect := "value not pointer or nil"
	got := GetAndSet(func() any { return 1 }, nil, time.Now().Add(time.Second))
	if got.Error() != expect {
		t.Errorf("expected %+v got %+v", expect, got)
	}

	var gotInt int
	expectInt := 1
	err := GetAndSet(func() any { return expectInt }, &gotInt, time.Now().Add(time.Second))
	if err != nil {
		t.Errorf("expected %+v got %+v", nil, err)
	}
	if err != nil {
		t.Errorf("expected %+v got %+v", expectInt, gotInt)
	}
}

func TestCache_Get(t *testing.T) {
	expectErrorMsg := "key not set"
	_, gotErr := Get("fake-key")
	if gotErr.Error() != expectErrorMsg {
		t.Errorf("expected %+v got %+v", expectErrorMsg, gotErr)
	}

	createdKey := "real-key"
	value := 1
	expirationTime := time.Now().Add(time.Second)
	_cache.records[createdKey] = record{value: value, created: true, expiresAt: &expirationTime}
	gotValue, gotErr := Get(createdKey)
	if gotErr != nil {
		t.Errorf("expected %+v got %+v", nil, gotErr)
	}
	if gotValue != value {
		t.Errorf("expected %+v got %+v", value, gotValue)
	}

	time.Sleep(time.Second)
	expectErrorMsg = "value expired"
	_, gotErr = Get(createdKey)
	if gotErr.Error() != expectErrorMsg {
		t.Errorf("expected %+v got %+v", expectErrorMsg, gotErr)
	}
}

func TestCache_Set(t *testing.T) {
	key := "key"
	Set(key, 1, nil)
	if _cache.records[key].value != 1 {
		t.Errorf("expected %+v got %+v", 1, _cache.records[key].value)
	}
	if _cache.records[key].created != true {
		t.Errorf("expected %+v got %+v", true, _cache.records[key].created)
	}
}

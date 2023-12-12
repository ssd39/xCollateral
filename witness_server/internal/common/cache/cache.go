package cache

import (
	"errors"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type record struct {
	value     any
	created   bool
	expiresAt *time.Time
}

type cache struct {
	records map[string]record
}

var _cache = &cache{records: map[string]record{}}
var mutex = sync.RWMutex{}

func GetAndSet(f func() any, v any, expirationTime time.Time) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("value not pointer or nil")
	}

	key := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	mutex.RLock()
	elem := _cache.records[key]
	mutex.RUnlock()
	if elem.created && (elem.expiresAt == nil || elem.expiresAt.After(time.Now())) {
		rv.Elem().Set(reflect.ValueOf(elem.value))
		log.Debug().Msgf("Cache hit for key %s with value %+v stored value %+v", key, v, elem.value)
	} else {
		newValue := f()
		if newValue != nil {
			mutex.Lock()
			_cache.records[key] = record{value: newValue, created: true, expiresAt: &expirationTime}
			mutex.Unlock()
			rv.Elem().Set(reflect.ValueOf(newValue))
			log.Debug().Msgf("Cache miss for key %s new value is %+v stored value %+v", key, v, newValue)
		} else {
			log.Debug().Msgf("Cache miss for key %s new value is %+v stored value %+v", key, v, elem.value)
		}
	}
	return nil
}

func Get(key string) (any, error) {
	mutex.RLock()
	record := _cache.records[key]
	mutex.RUnlock()
	if !record.created {
		return nil, errors.New("key not set")
	}
	if record.expiresAt != nil && record.expiresAt.Before(time.Now()) {
		return nil, errors.New("value expired")
	}

	return record.value, nil
}

func Set(key string, value any, expirationTime *time.Time) {
	mutex.Lock()
	_cache.records[key] = record{value: value, created: true, expiresAt: expirationTime}
	mutex.Unlock()
}

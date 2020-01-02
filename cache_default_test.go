package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	Key string
	Value string
}

func TestDefaultCacheSetSuccessWithString(t *testing.T) {
	key := "cache_key"
	val := "value"
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(string)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheSetSuccessWithInt(t *testing.T) {
	key := "cache_key"
	val := 1
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(int)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheSetSuccessWithBoolean(t *testing.T) {
	key := "cache_key"
	val := true
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(bool)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheSetSuccessWithStruct_set(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(testItem)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}


func TestDefaultCacheAddSuccessWithString(t *testing.T) {
	key := "cache_key"
	val := "value"
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(string)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheAddSuccessWithInt(t *testing.T) {
	key := "cache_key"
	val := 1
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(int)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheAddSuccessWithBoolean(t *testing.T) {
	key := "cache_key"
	val := true
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(bool)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheAddSuccessWithStruct(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	value, err := cache.Get(key)
	assert.NoError(t, err)

	cacheValue := new(testItem)
	err = json.Unmarshal(value, cacheValue)
	assert.NoError(t, err)
	assert.Equal(t, val, *cacheValue)
}

func TestDefaultCacheAddErrorCacheAlreadyExists(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.Error(t, err)
}

func TestDefaultCachePullSuccessWithStruct(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	_, err = cache.Pull(key)
	assert.NoError(t, err)
	assert.False(t, cache.Has(key))
}

func TestDefaultCacheExpired(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewDefaultCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	time.Sleep(5 * time.Second)
	_, err = cache.Pull(key)
	assert.Error(t, err)
}

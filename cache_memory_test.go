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

func TestMemoryCacheSetSuccessWithString(t *testing.T) {
	key := "cache_key"
	val := "value"
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheSetSuccessWithInt(t *testing.T) {
	key := "cache_key"
	val := 1
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheSetSuccessWithBoolean(t *testing.T) {
	key := "cache_key"
	val := true
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheSetSuccessWithStruct_set(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewMemoryCache(5 * time.Second)
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


func TestMemoryCacheAddSuccessWithString(t *testing.T) {
	key := "cache_key"
	val := "value"
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheAddSuccessWithInt(t *testing.T) {
	key := "cache_key"
	val := 1
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheAddSuccessWithBoolean(t *testing.T) {
	key := "cache_key"
	val := true
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheAddSuccessWithStruct(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewMemoryCache(5 * time.Second)
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

func TestMemoryCacheAddErrorCacheAlreadyExists(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewMemoryCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.Error(t, err)
}

func TestMemoryCachePullSuccessWithStruct(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewMemoryCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	_, err = cache.Pull(key)
	assert.NoError(t, err)
	assert.False(t, cache.Has(key))
}

func TestMemoryCacheExpired(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewMemoryCache(5 * time.Second)
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	time.Sleep(5 * time.Second)
	_, err = cache.Pull(key)
	assert.Error(t, err)
}

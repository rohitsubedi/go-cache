package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileCacheSetSuccessWithString(t *testing.T) {
	key := "cache_key"
	val := "value"
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheSetSuccessWithInt(t *testing.T) {
	key := "cache_key"
	val := 1
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheSetSuccessWithBoolean(t *testing.T) {
	key := "cache_key"
	val := true
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheSetSuccessWithStruct_set(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewFileCache(5 * time.Second, "cache")
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



func TestFileCacheAddSuccessWithString(t *testing.T) {
	key := "cache_key1"
	val := "value"
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheAddSuccessWithInt(t *testing.T) {
	key := "cache_key2"
	val := 1
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheAddSuccessWithBoolean(t *testing.T) {
	key := "cache_key3"
	val := true
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheAddSuccessWithStruct(t *testing.T) {
	key := "cache_key4"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewFileCache(5 * time.Second, "cache")
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

func TestFileCacheAddErrorCacheAlreadyExists(t *testing.T) {
	key := "cache_key5"
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

func TestFileCachePullSuccessWithStruct(t *testing.T) {
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

	_, err = cache.Pull(key)
	assert.NoError(t, err)
	assert.False(t, cache.Has(key))
}

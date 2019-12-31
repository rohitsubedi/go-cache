package cache

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisCacheSetSuccessWithString(t *testing.T) {
	key := "cache_key"
	val := "value"
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheSetSuccessWithInt(t *testing.T) {
	key := "cache_key"
	val := 1
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheSetSuccessWithBoolean(t *testing.T) {
	key := "cache_key"
	val := true
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheSetSuccessWithStruct_set(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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



func TestRedisCacheAddSuccessWithString(t *testing.T) {
	key := "cache_key1"
	val := "value"
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheAddSuccessWithInt(t *testing.T) {
	key := "cache_key2"
	val := 1
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheAddSuccessWithBoolean(t *testing.T) {
	key := "cache_key3"
	val := true
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheAddSuccessWithStruct(t *testing.T) {
	key := "cache_key4"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
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

func TestRedisCacheAddErrorCacheAlreadyExists(t *testing.T) {
	key := "cache_key5"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.NoError(t, err)

	err = cache.Add(key, val)
	assert.Error(t, err)
}

func TestRedisCachePullSuccessWithStruct(t *testing.T) {
	key := "cache_key"
	val := testItem{
		Key:   "Rohit",
		Value: "Subedi",
	}
	cache, err := NewRedisCache(5 * time.Second, "0.0.0.0:6379", "redis_password")
	assert.NoError(t, err)

	err = cache.Set(key, val)
	assert.NoError(t, err)
	assert.True(t, cache.Has(key))

	_, err = cache.Pull(key)
	assert.NoError(t, err)
	assert.False(t, cache.Has(key))
}

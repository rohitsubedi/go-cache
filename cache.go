package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-redis/redis/v7"
)

var (
	cacheTypeDefault  = "memory"
	cacheTypeFile     = "file"
	cacheTypeRedis    = "redis"
	cacheTypeMemcache = "memcache"
	defaultExpiration = 0 * time.Second
)

var (
	ErrCacheNotFound      = errors.New("cache lib: cache not found")
	ErrCacheExpired       = errors.New("cache lib: cache expired")
	ErrConnectingRedis    = errors.New("cache lib: cannot connect to redis server")
	ErrCreatingFile       = errors.New("cache lib: cannot create file on the given path")
	ErrCacheAlreadyExists = errors.New("cache lib: cache already exists")
)

type Cache interface {
	Add(key string, value interface{}) error
	Set(key string, value interface{}) error
	Get(key string) ([]byte, error)
	Pull(key string) ([]byte, error)
	Has(key string) bool
	Delete(key string)
	Flush()
}

type cacheItem struct {
	value      []byte
	expiration int64
}

type cache struct {
	mu             sync.RWMutex
	cacheType      string
	expiration     time.Duration
	items          map[string]cacheItem
	filePath       string
	cacheFiles     map[string]struct{}
	redisClient    *redis.Client
	memCacheClient *memcache.Client
	cleaner        *cacheCleaner
}

type cacheCleaner struct {
	interval *time.Timer
	stop     chan bool
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
func NewDefaultCache(expiration time.Duration) (Cache, error) {
	var cleaner *cacheCleaner

	if expiration <= defaultExpiration {
		expiration = defaultExpiration
	} else {
		cleaner = &cacheCleaner{
			interval: time.NewTimer(expiration),
			stop:     make(chan bool),
		}
	}

	cache := &cache{
		cacheType:  cacheTypeDefault,
		expiration: expiration,
		items:      make(map[string]cacheItem),
		cleaner:    cleaner,
	}

	cache.cleanExpiredCache()

	return cache, nil
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
// path string directory path where the cache file can be stored. It should have write permission
func NewFileCache(expiration time.Duration, path string) (Cache, error) {
	var cleaner *cacheCleaner

	if expiration <= defaultExpiration {
		expiration = defaultExpiration
	} else {
		cleaner = &cacheCleaner{
			interval: time.NewTimer(expiration),
			stop:     make(chan bool),
		}
	}

	cache := &cache{
		cacheType:  cacheTypeFile,
		expiration: expiration,
		filePath:   path,
		cacheFiles: make(map[string]struct{}),
		cleaner:    cleaner,
	}

	cache.cleanExpiredCache()

	return cache, nil
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
func NewRedisCache(expiration time.Duration, host, password string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("%v: %w", ErrConnectingRedis, err)
	}

	if expiration <= defaultExpiration {
		expiration = defaultExpiration
	}

	return &cache{
		cacheType:   cacheTypeRedis,
		expiration:  expiration,
		redisClient: client,
	}, nil
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
func NewMemCache(expiration time.Duration, server ...string) (Cache, error) {
	if expiration <= defaultExpiration {
		expiration = defaultExpiration
	}

	memCacheClient := memcache.New(server...)
	if err := memCacheClient.Ping(); err != nil {
		return nil, err
	}

	return &cache{
		cacheType:      cacheTypeMemcache,
		expiration:     expiration,
		memCacheClient: memCacheClient,
	}, nil
}

// Delete deletes cache for the given key
func (c *cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.cacheType {
	case cacheTypeDefault:
		delete(c.items, key)
	case cacheTypeFile:
		_ = os.Remove(c.filePath + "/" + key)
	case cacheTypeRedis:
		c.redisClient.Del(key)
	case cacheTypeMemcache:
		_ = c.memCacheClient.Delete(key)
	}
}

// Flush deletes all the existing cache
func (c *cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.cacheType {
	case cacheTypeDefault:
		c.items = make(map[string]cacheItem)
	case cacheTypeFile:
		for key := range c.cacheFiles {
			_ = os.Remove(c.filePath + "/" + key)
		}
	case cacheTypeRedis:
		c.redisClient.FlushAll()
	case cacheTypeMemcache:
		_ = c.memCacheClient.FlushAll()
	}
}

// This will set the value to the key depending on the cache type user selects (memory, file, redis).
// If cache already exists for given key, it will return error. Returns error if there are any
func (c *cache) Add(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.has(key) {
		return ErrCacheAlreadyExists
	}

	return c.set(key, value)
}

// This will set the value to the key depending on the cache type user selects (memory, file, redis).
// This will override the existing value in the cache. Returns error if there are any
func (c *cache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.set(key, value)
}

func (c *cache) set(key string, value interface{}) error {
	var expiration int64
	val, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return err
	}

	if c.expiration > defaultExpiration {
		expiration = time.Now().Add(c.expiration).UnixNano()
	}

	switch c.cacheType {
	case cacheTypeDefault:
		c.items[key] = cacheItem{
			value:      val,
			expiration: expiration,
		}
	case cacheTypeFile:
		file, err := os.OpenFile(c.filePath+"/"+key, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0644))
		if err != nil {
			return fmt.Errorf("%v: %w", ErrCreatingFile, err)
		}

		if _, err := file.Write(val); err != nil {
			return err
		}

		c.cacheFiles[key] = struct{}{}
	case cacheTypeRedis:
		if err := c.redisClient.Set(key, val, c.expiration).Err(); err != nil {
			return err
		}
	case cacheTypeMemcache:
		if err := c.memCacheClient.Set(&memcache.Item{
			Key:        key,
			Value:      val,
			Expiration: int32(c.expiration.Seconds()),
		}); err != nil {
			return err
		}
	}

	return nil
}

// This will return boolean if the cache exists and is valid
func (c *cache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.has(key)
}

func (c *cache) has(key string) bool {
	switch c.cacheType {
	case cacheTypeDefault:
		item, found := c.items[key]
		if !found {
			return false
		}

		if item.expiration > 0 {
			if time.Now().UnixNano() > item.expiration {
				delete(c.items, key)
				return false
			}
		}
	case cacheTypeFile:
		fileInfo, err := os.Stat(c.filePath + "/" + key)
		if err != nil {
			return false
		}

		if c.expiration > 0 && time.Now().UnixNano() > (c.expiration.Nanoseconds()+fileInfo.ModTime().UnixNano()) {
			_ = os.Remove(c.filePath + "/" + key)
			return false
		}
	case cacheTypeRedis:
		if _, err := c.redisClient.Get(key).Result(); err != nil {
			return false
		}
	case cacheTypeMemcache:
		if _, err := c.memCacheClient.Get(key); err != nil {
			return false
		}
	default:
		return false
	}

	return true
}

// This returns the value in the cache for the given key if its valid. Returns error if cache doesn'interval exist or expired
func (c *cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch c.cacheType {
	case cacheTypeDefault:
		return c.getDefaultCache(key, false)
	case cacheTypeFile:
		return c.getFileCache(key, false)
	case cacheTypeRedis:
		return c.getRedisCache(key, false)
	case cacheTypeMemcache:
		return c.getMemCache(key, false)
	}

	return nil, ErrCacheNotFound
}

// This returns the value in the cache for the given key if it's valid (AND also removes the cache for the given key).
// Returns error if cache doesn'interval exist or expired
func (c *cache) Pull(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch c.cacheType {
	case cacheTypeDefault:
		return c.getDefaultCache(key, true)
	case cacheTypeFile:
		return c.getFileCache(key, true)
	case cacheTypeRedis:
		return c.getRedisCache(key, true)
	case cacheTypeMemcache:
		return c.getMemCache(key, true)
	}

	return nil, ErrCacheNotFound
}

// Returns value from memory cache for given key. Removes current cache depending on second parameter
func (c *cache) getDefaultCache(key string, removeCurrent bool) ([]byte, error) {
	item, found := c.items[key]
	if !found {
		return nil, ErrCacheNotFound
	}

	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
			delete(c.items, key)
			return nil, ErrCacheExpired
		}
	}

	if removeCurrent {
		delete(c.items, key)
	}

	return item.value, nil
}

// Returns value from file cache for given key. Removes current cache depending on second parameter
func (c *cache) getFileCache(key string, removeCurrent bool) ([]byte, error) {
	fileInfo, err := os.Stat(c.filePath + "/" + key)
	if err != nil {
		return nil, ErrCacheNotFound
	}

	if c.expiration > 0 && time.Now().UnixNano() > (c.expiration.Nanoseconds()+fileInfo.ModTime().UnixNano()) {
		_ = os.Remove(c.filePath + "/" + key)
		return nil, ErrCacheExpired
	}

	value, err := ioutil.ReadFile(c.filePath + "/" + key)
	if err != nil {
		return nil, ErrCacheNotFound
	}

	if removeCurrent {
		_ = os.Remove(c.filePath + "/" + key)
	}

	return value, nil
}

// Returns value from redis cache for given key. Removes current cache depending on second parameter
func (c *cache) getRedisCache(key string, removeCurrent bool) ([]byte, error) {
	val, err := c.redisClient.Get(key).Result()
	if err != nil {
		return nil, ErrCacheNotFound
	}

	if removeCurrent {
		_ = c.redisClient.Del(key)
	}

	return []byte(val), nil
}

// Returns value from redis cache for given key. Removes current cache depending on second parameter
func (c *cache) getMemCache(key string, removeCurrent bool) ([]byte, error) {
	val, err := c.memCacheClient.Get(key)
	if err != nil {
		return nil, ErrCacheNotFound
	}

	if removeCurrent {
		_ = c.memCacheClient.Delete(key)
	}

	return val.Value, nil
}

// This is a job that will execute each duration of the cache and clears the expired cache
func (c *cache) cleanExpiredCache() {
	if c.cleaner == nil {
		return
	}

	if c.cacheType == cacheTypeMemcache || c.cacheType == cacheTypeRedis {
		return
	}

	runtime.SetFinalizer(c.cleaner, stopCleaningRoutine)

	go func() {
		for {
			select {
			case <-c.cleaner.interval.C:
				switch c.cacheType {
				case cacheTypeDefault:
					for key, _ := range c.items {
						c.has(key)
					}
				case cacheTypeFile:
					for key, _ := range c.cacheFiles {
						go c.has(key)
					}
				}

				c.cleaner.interval.Reset(c.expiration)
			case <-c.cleaner.stop:
				c.cleaner.interval.Stop()
			}

		}
	}()
}

// go routine is stopped stop is set to true
func stopCleaningRoutine(cleaner *cacheCleaner) {
	cleaner.stop <- true
}

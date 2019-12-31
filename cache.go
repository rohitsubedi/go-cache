package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
)

var (
	cacheTypeDefault  = "memory"
	cacheTypeFile     = "file"
	cacheTypeRedis    = "redis"
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
}

type item struct {
	value      []byte
	expiration int64
}

type cache struct {
	mu          sync.RWMutex
	cacheType   string
	expiration  time.Duration
	items       map[string]item
	filePath    string
	redisClient *redis.Client
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
func NewMemoryCache(expiration time.Duration) (Cache, error) {
	return &cache{
		cacheType:  cacheTypeDefault,
		expiration: expiration,
		items:      make(map[string]item),
	}, nil
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
// path string directory path where the cache file can be stored. It should have write permission
func NewFileCache(expiration time.Duration, path string) (Cache, error) {
	return &cache{
		cacheType:  cacheTypeFile,
		expiration: expiration,
		filePath:   path,
	}, nil
}

// expiration time.Duration duration for cache to expire. 0*time.Second indicates the cache will never expire
// path string directory path where the cache file can be stored. It should have write permission
func NewRedisCache(expiration time.Duration, host, password string) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("%v: %w", ErrConnectingRedis, err)
	}

	return &cache{
		cacheType:   cacheTypeRedis,
		expiration:  expiration,
		redisClient: client,
	}, nil
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
		c.items[key] = item{
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
	case cacheTypeRedis:
		if err := c.redisClient.Set(key, val, c.expiration).Err(); err != nil {
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
	default:
		return false
	}

	return true
}

// This returns the value in the cache for the given key if its valid. Returns error if cache doesn't exist or expired
func (c *cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch c.cacheType {
	case cacheTypeDefault:
		return c.getMemoryCache(key, false)
	case cacheTypeFile:
		return c.getFileCache(key, false)
	case cacheTypeRedis:
		return c.getRedisCache(key, false)
	}

	return nil, ErrCacheNotFound
}

// This returns the value in the cache for the given key if it's valid (AND also removes the cache for the given key).
// Returns error if cache doesn't exist or expired
func (c *cache) Pull(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	switch c.cacheType {
	case cacheTypeDefault:
		return c.getMemoryCache(key, true)
	case cacheTypeFile:
		return c.getFileCache(key, true)
	case cacheTypeRedis:
		return c.getRedisCache(key, true)
	}

	return nil, ErrCacheNotFound
}

// Returns value from memory cache for given key. Removes current cache depending on second parameter
func (c *cache) getMemoryCache(key string, removeCurrent bool) ([]byte, error) {
	item, found := c.items[key]
	if !found {
		return nil, ErrCacheNotFound
	}

	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
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

//type v struct {
//	Key   string `json:"key"`
//	Value string `json:"value"`
//}
//
//func main() {
//	cache, err := NewRedisCache(5*time.Second, "0.0.0.0:6379", "redis_password")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	value := v{
//		Key:   "Rohit",
//		Value: "Subedi.",
//	}
//
//	if err := cache.Set("rohit", value); err != nil {
//		fmt.Println(err)
//	}
//
//	cacheValue, err := cache.Get("rohit")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	original := new(v)
//	err = json.Unmarshal(cacheValue, original)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(original)
//
//	cacheValue, err = cache.Pull("rohit")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}

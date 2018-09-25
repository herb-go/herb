package cache

import "time"

type Cacheable interface {
	Set(key string, v interface{}, ttl time.Duration) error
	Get(key string, v interface{}) error
	SetBytesValue(key string, bytes []byte, ttl time.Duration) error
	GetBytesValue(key string) ([]byte, error)
	Del(key string) error
	IncrCounter(key string, increment int64, ttl time.Duration) (int64, error)
	SetCounter(key string, v int64, ttl time.Duration) error
	GetCounter(key string) (int64, error)
	DelCounter(key string) error
	Expire(key string, ttl time.Duration) error
	ExpireCounter(key string, ttl time.Duration) error
	Load(key string, v interface{}, ttl time.Duration, loader func() (interface{}, error)) error
	MGetBytesValue(keys ...string) (map[string][]byte, error)
	MSetBytesValue(map[string][]byte, time.Duration) error
	FinalKey(string) (string, error)
	Flush() error
	DefualtTTL() time.Duration
}

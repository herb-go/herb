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
	Load(key string, v interface{}, ttl time.Duration, closure func(v interface{}) error) error
	Flush() error
	DefualtTTL() time.Duration
}

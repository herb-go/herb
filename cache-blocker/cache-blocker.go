package blocker

import (
	"net"
	"net/http"
	"time"

	"strconv"

	"github.com/herb-go/herb/cache"
)

//StatusAny stand for any status
const StatusAny = 0
const defaultBlockedStatus = http.StatusTooManyRequests

func New(cache cache.Cacheable, Identifier func(r *http.Request) (string, error)) *Blocker {
	return &Blocker{
		config:        map[int]statusConfig{},
		Cache:         cache,
		StatusBlocked: defaultBlockedStatus,
		Identifier:    Identifier,
	}
}

type statusConfig struct {
	ttlSecond      int64
	max            int64
	cacheKeyPrefix string
}
type Blocker struct {
	config        map[int]statusConfig
	Cache         cache.Cacheable
	StatusBlocked int
	Identifier    func(r *http.Request) (string, error)
}

func (b *Blocker) Block(status int, max int64, ttl time.Duration) {
	ttlSecond := int64(ttl / time.Second)
	b.config[status] = statusConfig{
		max:            max,
		ttlSecond:      ttlSecond,
		cacheKeyPrefix: strconv.Itoa(status) + cache.KeyPrefix + strconv.FormatInt(ttlSecond, 10) + cache.KeyPrefix,
	}
}
func (b *Blocker) buildCacheKey(id string, status int, config statusConfig) string {
	timeHash := int64(time.Now().Unix() / config.ttlSecond)
	return config.cacheKeyPrefix + cache.KeyPrefix + id + cache.KeyPrefix + strconv.FormatInt(timeHash, 10)
}
func (b *Blocker) isBlocked(id string) bool {
	for k := range b.config {
		config, ok := b.config[k]
		if ok == true {
			key := b.buildCacheKey(id, k, config)
			count, err := b.Cache.GetCounter(key)
			if err != cache.ErrNotFound {
				if err != nil {
					panic(err)
				}
				if count >= config.max {
					return true
				}
			}
		}
	}
	return false
}
func (b *Blocker) incr(ip string, status int) {
	checklist := []int{status, StatusAny}
	for k := range checklist {
		config, ok := b.config[checklist[k]]
		if ok == true {
			key := b.buildCacheKey(ip, status, config)
			_, err := b.Cache.IncrCounter(key, 1, time.Duration(config.ttlSecond)*time.Second)
			if err != nil {
				panic(err)
			}
		}
	}
}
func IPIdentifier(r *http.Request) (string, error) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip, nil
}
func (b *Blocker) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := b.Identifier(r)
	if err != nil {
		panic(err)
	}
	if b.isBlocked(id) {
		http.Error(w, http.StatusText(b.StatusBlocked), b.StatusBlocked)
		return
	}
	writer := blockWriter{
		w,
		200,
	}
	next(&writer, r)
	b.incr(id, writer.status)
}

type blockWriter struct {
	http.ResponseWriter
	status int
}

func (w *blockWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

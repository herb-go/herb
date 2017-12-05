package ipblocker

import (
	"net"
	"net/http"
	"time"

	"strconv"

	"github.com/herb-go/herb/cache"
)

const StatusAny = 0
const defaultBlockedStatus = http.StatusTooManyRequests

func New(name string, cache cache.Cacheable) *Blocker {
	return &Blocker{
		config:        map[int]statusConfig{},
		Cache:         cache,
		StatusBlocked: defaultBlockedStatus,
		Name:          name,
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
	Name          string
}

func (b *Blocker) Flush() error {
	return b.Cache.Flush()
}
func (b *Blocker) Block(status int, max int64, ttl time.Duration) {
	ttlSecond := int64(ttl / time.Second)
	b.config[status] = statusConfig{
		max:            max,
		ttlSecond:      ttlSecond,
		cacheKeyPrefix: b.Name + "-" + strconv.Itoa(status) + "-" + strconv.FormatInt(ttlSecond, 10) + "-",
	}
}
func (b *Blocker) buildCacheKey(ip string, status int, config statusConfig) string {
	timeHash := int64(time.Now().Unix() / config.ttlSecond)
	return config.cacheKeyPrefix + ip + "-" + strconv.FormatInt(timeHash, 10)
}
func (b *Blocker) isIpBlocked(ip string) bool {
	for k := range b.config {
		config, ok := b.config[k]
		if ok == true {
			key := b.buildCacheKey(ip, k, config)
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
func (b *Blocker) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if b.isIpBlocked(ip) {
		http.Error(w, http.StatusText(b.StatusBlocked), b.StatusBlocked)
		return
	}
	writer := blockWriter{
		w,
		200,
	}
	next(&writer, r)
	b.incr(ip, writer.status)
}

type blockWriter struct {
	http.ResponseWriter
	status int
}

func (w *blockWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

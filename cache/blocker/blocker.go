package blocker

import (
	"net"
	"net/http"
	"time"

	"strconv"

	"github.com/herb-go/herb/cache"
)

//StatusAny stand for any status when block.
const StatusAny = 0

//StatusAnyError stand for any status greater than 400 when block.
const StatusAnyError = -1
const defaultBlockedStatus = http.StatusTooManyRequests

//New create blocker with given cache and http request udentifier
func New(cache cache.Cacheable) *Blocker {
	return &Blocker{
		config:            map[int]statusConfig{},
		Cache:             cache,
		StatusCodeBlocked: defaultBlockedStatus,
		Identifier:        IPIdentifier,
	}
}

type statusConfig struct {
	ttlSecond      int64
	max            int64
	cacheKeyPrefix string
}

//Blocker blocker struct.
type Blocker struct {
	config map[int]statusConfig
	//Cache cache which store blcok data
	Cache cache.Cacheable
	//StatusCodeBlocked error status which will returned when request blcoker.Default value is 429.
	StatusCodeBlocked int
	//Identifier http request identifier
	Identifier func(r *http.Request) (string, error)
	//OnBlock acitons execed when access blocked
	OnBlock func(w http.ResponseWriter, r *http.Request)
}

//Block block config method.
//Requester request for morethan param max request which response staus is param status in param ttl will be blocked.
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

//DefaultBlockAction default block
func (b *Blocker) DefaultBlockAction(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(b.StatusCodeBlocked), b.StatusCodeBlocked)
}
func (b *Blocker) incr(ip string, status int) {
	checklist := []int{status, StatusAny}
	if status >= 400 {
		checklist = append(checklist, StatusAnyError)
	}
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

//IPIdentifier identify http request by ip address.
func IPIdentifier(r *http.Request) (string, error) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip, nil
}

//ServeMiddleware serve blocker as a middleware.
func (b *Blocker) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := b.Identifier(r)
	if err != nil {
		panic(err)
	}
	if b.isBlocked(id) {
		if b.OnBlock != nil {
			b.OnBlock(w, r)
		} else {
			b.DefaultBlockAction(w, r)
		}
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

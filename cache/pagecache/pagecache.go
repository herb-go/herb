package pagecache

import (
	"bytes"
	"net/http"
	"strings"

	"time"

	"github.com/herb-go/herb/cache"
)

var Debug bool

var PageCacheKeyHeader = "herbgo-debug-pagecache"

func New(c cache.Cacheable) *PageCache {
	p := PageCache{
		Cache:           c,
		KeyPrefix:       defaultKeyPrefix,
		StatusValidator: defualtStatusValidator,
	}
	return &p
}

type cachedPage struct {
	Status   int
	Header   map[string][]string
	Response []byte
}

type PageCache struct {
	Cache           cache.Cacheable
	KeyPrefix       string
	StatusValidator func(status int) bool
}

var defaultKeyPrefix = ""

func defualtStatusValidator(status int) bool {
	return status < 500
}
func (p *PageCache) ValidateStatus(status int) bool {
	if p.StatusValidator != nil {
		return p.StatusValidator(status)
	}
	return defualtStatusValidator(status)
}
func (p *PageCache) serve(key string, ttl time.Duration, w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	page := cachedPage{}
	err := p.Cache.Load(key, &page, ttl, func() (interface{}, error) {
		cw := cacheResponseWriter{
			writer: *(bytes.NewBuffer([]byte{})),
			header: http.Header{},
			status: 0,
		}
		next(&cw, r)
		page.Header = map[string][]string(cw.Header())
		page.Response = cw.writer.Bytes()
		page.Status = cw.status
		if p.ValidateStatus(cw.status) {
			return page, nil
		}
		return nil, cache.ErrNotCacheable
	})
	if err != nil {
		if err != cache.ErrEntryTooLarge && err != cache.ErrNotCacheable {
			panic(err)
		}
	}
	h := w.Header()
	for k, v := range page.Header {
		h[k] = v
	}
	if Debug {
		k, err := p.Cache.FinalKey(key)
		if err == nil {
			h.Set(PageCacheKeyHeader, strings.Replace(k, cache.KeyPrefix, " ", -1))
		}
	}
	if page.Status != 0 {
		w.WriteHeader(page.Status)
	}
	_, err = w.Write(page.Response)
	if err != nil {
		panic(err)
	}
}
func (p *PageCache) Middleware(keyGenerator func(r *http.Request) string, ttl time.Duration) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		key := keyGenerator(r)
		if key == "" {
			next(w, r)
			return
		}
		p.serve(p.KeyPrefix+key, ttl, w, r, next)
	}
}
func FieldMiddleware(FieldGenerator func(r *http.Request) *cache.Field, ttl time.Duration, statusValidator func(status int) bool) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		f := FieldGenerator(r)
		if f == nil {
			next(w, r)
			return
		}
		p := &PageCache{
			Cache: f.Cache,
		}
		p.serve(f.FieldName, ttl, w, r, next)
	}
}

type cacheResponseWriter struct {
	writer bytes.Buffer
	header http.Header
	status int
}

func (w *cacheResponseWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}
func (w *cacheResponseWriter) Header() http.Header {
	return w.header
}
func (w *cacheResponseWriter) WriteHeader(status int) {
	w.status = status
}

package cache_test

import (
	"testing"

	"github.com/herb-go/herb/cache"
)

func NewEntiyTooLargeCache(ttl int64) cache.Cacheable {
	config := cache.ConfigMap{}
	config.Set("Size", 1)
	c := cache.New()
	oc := &cache.OptionConfigMap{
		Driver:    "syncmapcache",
		TTL:       int64(ttl),
		Config:    config,
		Marshaler: "json",
	}
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c

}
func newTestCache(ttl int64) cache.Cacheable {
	config := &cache.ConfigMap{}
	config.Set("Size", 10000000)
	c := cache.New()
	oc := &cache.OptionConfigMap{
		Driver:    "syncmapcache",
		TTL:       int64(ttl),
		Config:    nil,
		Marshaler: "json",
	}
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	return c
}

var testLoader = func(key string) (interface{}, error) {
	return key, nil
}

func TestCacheEmptyKey(t *testing.T) {
	var err error
	c := newTestCache(3600)
	var data = []byte{}
	err = c.Set("", data, 0)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.Update("", data, 0)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.Get("", data)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.SetBytesValue("", data, 0)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.UpdateBytesValue("", data, 0)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("")
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.Del("")
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.SetCounter("", 0, 0)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}

	_, err = c.GetCounter("")
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	_, err = c.IncrCounter("", 0, 0)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.DelCounter("")
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.Expire("", 1000)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.ExpireCounter("", 1000)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
	err = c.Load("", nil, 0, testLoader)
	if err != cache.ErrKeyUnavailable {
		t.Fatal(err)
	}
}
func TestLoader(t *testing.T) {
	var err error
	var result string
	c := newTestCache(3600)
	result = ""
	err = c.Load("test", &result, 0, testLoader)
	if err != nil {
		t.Fatal(err)
	}
	if result != "test" {
		t.Fatal(result)
	}
	err = c.Load("test", &result, 0, testLoader)
	if err != nil {
		t.Fatal(err)
	}
	if result != "test" {
		t.Fatal(result)
	}
}

func TestLoaderEntiyTooLarge(t *testing.T) {
	var err error
	var result string
	c := NewEntiyTooLargeCache(3600)
	result = ""
	err = c.Load("test", &result, 0, testLoader)
	if err != nil {
		t.Fatal(err)
	}
	if result != "test" {
		t.Fatal(result)
	}
	err = c.Load("test", &result, 0, testLoader)
	if err != nil {
		t.Fatal(err)
	}
	if result != "test" {
		t.Fatal(result)
	}
}

func TestNotFound(t *testing.T) {
	var err error
	c := newTestCache(3600)
	result := []byte{}
	err = c.Expire("notexists", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = c.ExpireCounter("notexists", 0)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Update("notexists", result, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Get("notexists", &result)
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	err = c.UpdateBytesValue("notexists", result, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("notexists")
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	err = c.Del("notexists")
	if err != nil {
		t.Fatal(err)
	}
	err = c.DelCounter("notexists")
	if err != nil {
		t.Fatal(err)
	}
}

func TestFinalKey(t *testing.T) {
	c := newTestCache(3600)
	k, err := c.FinalKey("key")
	if err != nil {
		t.Fatal(err)
	}
	if k != cache.KeyPrefix+"key" {
		t.Fatal(k)
	}
}

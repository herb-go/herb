package syncmapcache

import (
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
)

func newGCTestCache(ttl int64) *cache.Cache {
	config := cache.ConfigMap{}
	config.Set("Size", 10000000)
	config.Set("CleanupIntervalInSecond", 2)
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

func TestFlush(t *testing.T) {
	c := newGCTestCache(300)
	d := c.Driver.(*Cache)
	err := c.SetBytesValue("test", []byte("test"), 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("test")
	if err != nil {
		t.Fatal(err)
	}
	if d.used != 4 {
		t.Fatal(d.used)
	}
	err = c.Flush()
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("test")
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	if d.used != 0 {
		t.Fatal(d.used)
	}
}

func TestDel(t *testing.T) {
	c := newGCTestCache(300)
	d := c.Driver.(*Cache)
	err := c.SetBytesValue("test", []byte("test"), 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("test")
	if err != nil {
		t.Fatal(err)
	}
	if d.used != 4 {
		t.Fatal(d.used)
	}
	err = c.Del("test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("test")
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	if d.used != 0 {
		t.Fatal(d.used)
	}
}
func TestGc(t *testing.T) {
	c := newGCTestCache(300)
	d := c.Driver.(*Cache)
	err := c.SetBytesValue("test", []byte("test"), 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("test")
	if err != nil {
		t.Fatal(err)
	}
	if d.used != 4 {
		t.Fatal(d.used)
	}
	time.Sleep(100 * time.Microsecond)
	time.Sleep(2 * time.Second)
	_, err = c.GetBytesValue("test")
	if err != nil {
		t.Fatal(err)
	}
	if d.used != 4 {
		t.Fatal(d.used)
	}
	time.Sleep(1 * time.Second)
	_, err = c.GetBytesValue("test")
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	if d.used != 4 {
		t.Fatal(d.used)
	}
	time.Sleep(1 * time.Second)
	if d.used != 0 {
		t.Fatal(d.used)
	}
}

func newAutoRemoveTestCache(ttl int64) *cache.Cache {
	config := cache.ConfigMap{}
	config.Set("Size", 6)
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

func TestAutoRemove(t *testing.T) {
	c := newAutoRemoveTestCache(300)
	d := c.Driver.(*Cache)
	err := c.SetBytesValue("test", []byte("test"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if d.used != 4 {
		t.Fatal(d.used)
	}
	err = c.SetBytesValue("tes", []byte("tes"), 0)
	if err != nil {
		t.Fatal(err)
	}
	if d.used != 3 {
		t.Fatal(d.used)
	}
	_, err = c.GetBytesValue("tes")
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetBytesValue("test")
	if err != cache.ErrNotFound {
		t.Fatal(err)
	}
	err = c.SetBytesValue("testtest", []byte("testtest"), 0)
	if err != cache.ErrEntryTooLarge {
		t.Fatal(err)
	}

}

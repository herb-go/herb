package cache_test

import (
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/drivers/syncmapcache"
)

func TestDriver(t *testing.T) {
	factories := cache.Factories()
	if len(factories) != 2 {
		t.Fatal(factories)
	}
	dc, err := cache.NewDriver("dummycache", nil, "")
	if err != nil {
		t.Errorf("New driver  error %s", err)
	}
	if dc == nil {
		t.Fatal(dc)
	}
	cache.UnregisterAll()
	factories = cache.Factories()
	if len(factories) != 0 {
		t.Fatal(factories)
	}
	cache.Register("dummycache", func(conf cache.Config, prefix string) (cache.Driver, error) {
		return &cache.DummyCache{}, nil
	})
	syncmapcache.Register()
}

func TestEmptyDriver(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
	}()
	cache.Register("test", nil)

}
func TestDupdriver(t *testing.T) {
	var stage = 0
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
		if stage != 1 {
			t.Fatal(stage)
		}
	}()
	var testfactory = func(conf cache.Config, prefix string) (cache.Driver, error) {
		return nil, nil
	}
	cache.Register("test", testfactory)
	stage = 1
	cache.Register("test", testfactory)
	stage = 2
}

func TestNotExistDriver(t *testing.T) {
	_, err := cache.NewDriver("notexist", nil, "")
	if err == nil {
		t.Fatal(err)
	}
}

func TestMustNewDriver(t *testing.T) {
	var stage = 0
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
		if stage != 1 {
			t.Fatal(stage)
		}
	}()

	_ = cache.MustNewDriver("dummycache", nil, "")
	stage = 1
	_ = cache.MustNewDriver("notexist", nil, "")
	stage = 2
}

func TestNewSubCache(t *testing.T) {
	prefix := "test"
	c := &cache.ConfigJSON{}
	c.Set(prefix+"Driver", "syncmapcache")
	c.Set(prefix+"TTL", -1)
	c.Set(prefix+"Config.Size", 100000)
	c.Set(prefix+"Marshaler", "json")
	ca, err := cache.NewSubCache(c, prefix)
	if err != nil {
		panic(err)
	}
	if ca.TTL != -1*time.Second {
		t.Fatal(ca.TTL)
	}
	driver := ca.Driver.(*syncmapcache.Cache)
	if driver == nil {
		t.Fatal(driver)
	}
	if driver.Size != 100000 {
		t.Fatal(driver)
	}
}

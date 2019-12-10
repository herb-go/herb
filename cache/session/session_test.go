package session

import (
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
)

func NewTestCacheStore(ttl time.Duration) *Store {
	c := cache.New()
	oc := cache.NewOptionConfig()
	oc.Driver = "syncmapcache"
	oc.TTL = int64(ttl) * int64(time.Second)
	oc.Config = nil
	oc.Marshaler = "json"
	err := c.Init(oc)
	if err != nil {
		panic(err)
	}
	err = c.Flush()
	if err != nil {
		panic(err)
	}
	s := MustCacheStore(c, ttl)
	m, err := cache.NewMarshaler("json")
	if err != nil {
		panic(err)
	}
	s.Marshaler = m
	return s
}

func TestFlag(t *testing.T) {
	store := New()
	s := NewSession("", store)
	s.SetFlag(FlagTemporay, true)
	if s.HasFlag(FlagTemporay) == false {
		t.Fatal(s.HasFlag(FlagTemporay))
	}
	s.SetFlag(FlagTemporay, false)
	if s.HasFlag(FlagTemporay) == true {
		t.Fatal(s.HasFlag(FlagTemporay))
	}
}

func TestErrNilPointerAndRegenerate(t *testing.T) {
	var err error
	var result string
	store := NewTestCacheStore(1)
	s := NewSession("test", store)
	err = s.Set("testkey", "testvalue")
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = s.Get("testkey", nil)
	if err != ErrNilPointer {
		t.Fatal(err)
	}
	err = s.Get("testkey", result)
	if err != ErrNilPointer {
		t.Fatal(err)
	}
	err = s.Get("testkey", &result)
	if err != nil {
		t.Fatal(err)
	}

	if result != "testvalue" {
		t.Fatal(result)
	}
	err = s.Load()
	if err != nil {
		t.Fatal(err)
	}
	s.Regenerate()
	err = s.Get("testkey", &result)
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	token, err := s.Token()
	if err != nil {
		t.Fatal(err)
	}
	if token != "test" {
		t.Fatal(token)
	}
	time.Sleep(1001 * time.Millisecond)
	s = NewSession("test", store)

	err = s.Load()
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	err = s.Load()
	if err != ErrDataNotFound {
		t.Fatal(err)
	}
	err = s.Del("notexist")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEmptyKey(t *testing.T) {
	var err error
	store := NewTestCacheStore(1)
	s := NewSession("test", store)
	err = s.DeleteAndSave()
	if err != nil {
		t.Fatal(err)
	}
	token, err := s.Token()
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		t.Fatal(token)
	}
	err = s.Load()
	if err != ErrTokenNotValidated {
		t.Fatal(err)
	}
	result := ""
	err = s.Get("testkey", &result)
	if err != ErrTokenNotValidated {
		t.Fatal(err)
	}
	s = NewSession("test", store)
	s.token = ""
	err = store.LoadSession(s)
	if err != ErrTokenNotValidated {
		t.Fatal(err)
	}
}

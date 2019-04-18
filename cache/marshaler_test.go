package cache_test

import (
	"strings"
	"testing"

	"github.com/herb-go/herb/cache"
)

func TestMarshaler(t *testing.T) {
	marshalers := cache.MarshalerFactories()
	if len(marshalers) != 1 {
		t.Fatal(marshalers)
	}
	cache.UnregisterAllMarshalers()
	marshalers = cache.MarshalerFactories()
	if len(marshalers) != 0 {
		t.Fatal(marshalers)
	}
	cache.RegisterMarshaler("json", func() (cache.Marshaler, error) {
		return &cache.JSONMarshaler{}, nil
	})
}

func TestEmpteyTestMarshler(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
	}()
	cache.RegisterMarshaler("test", nil)

}
func TestDupTestMarshler(t *testing.T) {
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
	var testmarshaler = func() (cache.Marshaler, error) {
		return nil, nil
	}
	cache.RegisterMarshaler("test", testmarshaler)
	stage = 1
	cache.RegisterMarshaler("test", testmarshaler)
	stage = 2
}

func TestNotExistMarshler(t *testing.T) {
	_, err := cache.NewMarshaler("notexist")
	if err == nil {
		t.Fatal(err)
	}
	if strings.Contains(err.Error(), "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler") {
		t.Fatal(err)
	}
	_, err = cache.NewMarshaler("msgpack")
	if err == nil {
		t.Fatal(err)
	}
	if !strings.Contains(err.Error(), "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler") {
		t.Fatal(err)
	}
}

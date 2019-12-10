package cache_test

import (
	"strings"
	"testing"

	"github.com/herb-go/herb/cache"
)

func TestOption(t *testing.T) {
	oc := cache.NewOptionConfig()
	oc.Driver = "dummycache"
	oc.TTL = int64(3600)
	oc.Config = nil
	c := cache.New()
	err := c.Init(oc)
	if err == nil {
		t.Fatal(err)
	}
	if !strings.Contains(err.Error(), "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler") {
		t.Fatal(err)
	}
}

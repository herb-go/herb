package cache_test

import (
	"strings"
	"testing"

	"github.com/herb-go/herb/cache"
)

func TestOption(t *testing.T) {
	o := &cache.OptionConfigJSON{}
	o.Driver = "dummycache"
	o.Config = nil
	o.TTL = 3600
	c := cache.New()
	err := c.Init(o)
	if err == nil {
		t.Fatal(err)
	}
	if !strings.Contains(err.Error(), "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler") {
		t.Fatal(err)
	}
}

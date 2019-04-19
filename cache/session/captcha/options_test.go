package captcha

import (
	"testing"
	"time"

	"github.com/herb-go/herb/cache"

	_ "github.com/herb-go/herb/cache/drivers/syncmapcache"
	_ "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler"

	"github.com/herb-go/herb/cache/session"
)

func NewCatpcha() *Captcha {
	config := &cache.ConfigJSON{}
	config.Set("Size", 10000000)
	sc := cache.New()
	oc := &cache.OptionConfigMap{
		Driver:    "syncmapcache",
		TTL:       3600,
		Config:    nil,
		Marshaler: "json",
	}
	err := sc.Init(oc)
	if err != nil {
		panic(err)
	}
	s := session.MustCacheStore(sc, time.Hour)
	captcha := New(s)
	c := &Config{}
	c.Enabled = true
	c.Driver = "testcaptcha"
	c.ApplyTo(captcha)
	return captcha
}

func TestConfig(t *testing.T) {
	c := NewCatpcha()
	if c == nil {
		t.Fatal(c)
	}
	if c.Enabled == false {
		t.Fatal(c)
	}
}

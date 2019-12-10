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
	captcha := newEmptyCaptcha()
	c := &Config{}
	c.Enabled = true
	c.Driver = "testcaptcha"
	err := c.ApplyTo(captcha)
	if err != nil {
		panic(err)
	}
	return captcha
}
func newEmptyCaptcha() *Captcha {
	sc := cache.New()
	oc := cache.NewOptionConfig()
	oc.Driver = "syncmapcache"
	oc.TTL = 3600
	oc.Config = nil
	oc.Marshaler = "json"
	err := sc.Init(oc)
	if err != nil {
		panic(err)
	}
	s := session.MustCacheStore(sc, time.Hour)
	captcha := New(s)
	return captcha
}
func TestConfig(t *testing.T) {
	captcha := newEmptyCaptcha()
	c := &Config{
		Driver:         "testcaptcha",
		Enabled:        false,
		AddrWhiteList:  []string{"test"},
		DisabledScenes: map[string]bool{"test": false},
	}
	err := c.ApplyTo(captcha)
	if err != nil {
		t.Fatal(err)
	}
	if captcha.Enabled != false {
		t.Fatal(captcha)
	}
	if len(captcha.AddrWhiteList) != len(c.AddrWhiteList) {
		t.Fatal(captcha)
	}
	if len(captcha.DisabledScenes) != len(c.DisabledScenes) {
		t.Fatal(captcha)
	}
}

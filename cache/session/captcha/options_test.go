package captcha

import (
	"testing"
	"time"

	_ "github.com/herb-go/herb/cache/marshalers/msgpackmarshaler"
	"github.com/herb-go/herb/cache/session"
)

func NewCatpcha() *Captcha {
	s := session.MustClientStore([]byte("test"), time.Hour)
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

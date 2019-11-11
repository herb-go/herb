package apiserver

import (
	"testing"
	"time"

	"github.com/herb-go/herb/service/httpservice"
)

func TestConfig(t *testing.T) {
	var err error
	Reset()
	o := newOption()
	as := o.server()
	as.CleanConfig()
	defer func() {
		Reset()
		as.CleanConfig()
	}()
	if as.config != nil {
		t.Fatal(defaultConfig)
	}
	c := as.Config()
	if c != defaultConfig {
		t.Fatal(c)
	}
	as.CleanConfig()
	config := &httpservice.Config{}
	err = as.SetConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	c = as.Config()
	if c != config {
		t.Fatal(c)
	}
	err = as.SetConfig(config)
	if err == nil {
		t.Fatal(err)
	}
	as.Start("", nil)
	as.Stop("")
	time.Sleep(time.Millisecond)
	err = as.SetConfig(config)
	if err != nil {
		t.Fatal(err)
	}
}

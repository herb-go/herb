package captcha

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache/session"
)

type testDriver struct {
}

func (d *testDriver) Name() string {
	return "test"
}
func (d *testDriver) MustCaptcha(s *session.Store, w http.ResponseWriter, r *http.Request, scene string, reset bool) {
	var code string
	err := s.Get(r, "captcha", &code)
	if err == session.ErrDataNotFound {
		code = ""
		err = nil
	}
	if err != nil {
		panic(err)
	}
	if code == "" || reset {
		code = strconv.FormatInt(time.Now().Unix(), 10)
		err = s.Set(r, "captcha", code)
	}
	output, err := json.Marshal(map[string]interface{}{"Code": code})
	if err != nil {
		panic(err)
	}
	w.Write([]byte(output))
}
func (d *testDriver) Verify(s *session.Store, r *http.Request, scene string, token string) (bool, error) {
	var code string
	err := s.Get(r, "captcha", &code)
	if err == session.ErrDataNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return code == token, nil
}

func registerTestDriver() {
	Register("testcaptcha", func(conf cache.Config, prefix string) (Driver, error) {
		return &testDriver{}, nil
	})
}

func TestRegisterDriver(t *testing.T) {
	fs := Factories()
	if len(fs) != 1 {
		t.Fatal(fs)
	}
	UnregisterAll()
	fs = Factories()
	if len(fs) != 0 {
		t.Fatal(fs)
	}
	registerTestDriver()
}

func TestRegisterEmptyDriver(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
	}()
	Register("test", nil)
}

func TestRegisterDupDriver(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal(r)
		}
	}()
	Register("testcaptcha", func(conf cache.Config, prefix string) (Driver, error) {
		return &testDriver{}, nil
	})
}

func init() {
	registerTestDriver()
}

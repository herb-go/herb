package captcha

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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

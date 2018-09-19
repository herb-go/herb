package captcha

import (
	"net/http"
	"sync"

	"github.com/herb-go/herb/cache/session"
)

const HeaderReset = "X-Reset-Captcha"
const HeaderCaptchaName = "X-Captcha-Name"
const HeaderCaptchaEnabled = "X-Captcha-Enabled"

var (
	factorysMu sync.RWMutex
	factories  = make(map[string]Factory)
)

func defaultEnabledChecker(captcha *Captcha, scene string, w http.ResponseWriter, r *http.Request) (bool, error) {
	return true, nil
}

func New(s *session.Store) *Captcha {
	return &Captcha{
		DisabledScenes: map[string]bool{},
		Session:        s,
		EnabledChecker: defaultEnabledChecker,
	}
}

type Captcha struct {
	driver         Driver
	Session        *session.Store
	Enabled        bool
	DisabledScenes map[string]bool
	EnabledChecker func(captcha *Captcha, scene string, w http.ResponseWriter, r *http.Request) (bool, error)
}

func (c *Captcha) EnabledCheck(scene string, w http.ResponseWriter, r *http.Request) (bool, error) {
	if !c.Enabled || c.DisabledScenes[scene] {
		return false, nil
	}
	return c.EnabledChecker(c, scene, w, r)
}

func (c *Captcha) CaptchaAction(scene string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enabled, err := c.EnabledCheck(scene, w, r)
		if err != nil {
			panic(err)
		}
		w.Header().Set(HeaderCaptchaName, c.driver.Name())
		if enabled {
			w.Header().Set(HeaderCaptchaEnabled, "Enabled")
		}
		if enabled {
			c.driver.MustCaptcha(scene, r.Header.Get(HeaderReset) != "", w, r)
			return
		}
		_, err = w.Write([]byte{})
		if err != nil {
			panic(err)
		}
	}
}

func (c *Captcha) Verify(r *http.Request, scene string, token string) (bool, error) {
	return c.driver.Verify(r, scene, token)
}

func (c *Captcha) Verifier(r *http.Request, scene string) Verifier {
	return func(token string) (bool, error) {
		return c.Verify(r, scene, token)
	}
}

type Verifier func(token string) (bool, error)

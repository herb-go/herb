package captcha

import (
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/herb-go/herb/cache/session"
)

//HeaderReset header which should be passed in to reset captcha.
const HeaderReset = "X-Reset-Captcha"

//HeaderCaptchaName header which contains the captcha name which used.
const HeaderCaptchaName = "X-Captcha-Name"

//HeaderCaptchaEnabled header which contians if captcha is enabeld.
//If catpcha enabled,value "enabled" will be passed.
//Otherwise no header is passed.
const HeaderCaptchaEnabled = "X-Captcha-Enabled"

var (
	factorysMu sync.RWMutex
	factories  = make(map[string]Factory)
)

func defaultEnabledChecker(captcha *Captcha, scene string, r *http.Request) (bool, error) {
	return true, nil
}

//New create a new empty captcha instance with given session store.
func New(s *session.Store) *Captcha {
	return &Captcha{
		DisabledScenes: map[string]bool{},
		SessionStore:   s,
		EnabledChecker: defaultEnabledChecker,
		AddrWhiteList:  []string{},
	}
}

//Captcha  captcha struct.
type Captcha struct {
	driver Driver
	//Session captcha session store.
	SessionStore *session.Store
	//Enabled if captcha is enabled.
	Enabled bool
	//AddrWhiteList ip addr white list.Ip start with value in list doesn't need captcha.
	AddrWhiteList []string
	//DisabledScenes scenes which doesn't neec captcha.
	DisabledScenes map[string]bool
	//EnabledChecker function which check if captcha is necessarily.
	EnabledChecker func(captcha *Captcha, scene string, r *http.Request) (bool, error)
}

//EnabledCheck check if http request in given scene need captcha.
//Return true if captcha is necessarily,and any error if raised.
func (c *Captcha) EnabledCheck(scene string, r *http.Request) (bool, error) {
	if !c.Enabled || c.DisabledScenes[scene] {
		return false, nil
	}
	if len(c.AddrWhiteList) > 0 {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return false, err
		}
		for k := range c.AddrWhiteList {
			if strings.HasPrefix(host, c.AddrWhiteList[k]) {
				return false, nil
			}
		}
	}
	return c.EnabledChecker(c, scene, r)
}

//CaptchaAction action which afford capcha.
//Return captcha config json or empty object json if  doesn't need captcha.
//If reset header is passed in,captcha will be reseted if supported.
func (c *Captcha) CaptchaAction(scene string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		enabled, err := c.EnabledCheck(scene, r)
		if err != nil {
			panic(err)
		}
		w.Header().Set(HeaderCaptchaName, c.driver.Name())
		if enabled {
			w.Header().Set(HeaderCaptchaEnabled, "Enabled")
		}
		if enabled {
			c.driver.MustCaptcha(c.SessionStore, w, r, scene, r.Header.Get(HeaderReset) != "")
			return
		}
		_, err = w.Write([]byte("{}"))
		if err != nil {
			panic(err)
		}
	}
}

//Verify verify if token is validated with given http rquest and scene.
//return verify result and any error raised.
func (c *Captcha) Verify(r *http.Request, scene string, token string) (bool, error) {
	e, err := c.EnabledCheck(scene, r)
	if err != nil {
		return false, err
	}
	if !e {
		return true, nil
	}
	return c.driver.Verify(c.SessionStore, r, scene, token)
}

//Verifier create verifier with given http request and scene.
func (c *Captcha) Verifier(r *http.Request, scene string) Verifier {
	return func(token string) (bool, error) {
		return c.Verify(r, scene, token)
	}
}

//Verifier verifier interface.
type Verifier func(token string) (bool, error)

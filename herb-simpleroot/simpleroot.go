package simpleroot

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"time"

	"sync"

	"github.com/herb-go/herb/cache"
)

const DefaultRemindDuration = 5 * 24 * time.Hour
const DefaultInactiveDuration = -1
const DefaultLoginFailStatus = http.StatusUnprocessableEntity
const DefaultCookieName = "herb-simpleroot-token"
const CacheTokenKey = "Token"
const CacheRefreshedTimeKey = "Refreshed"

var Tokenlength = 64
var TokenMask = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var tokenLock = sync.Mutex{}

type UserVerifier interface {
	Verify(r *http.Request) (Verified bool, err error)
	Install(r *http.Request) (success bool, err error)
}

func New(verifier UserVerifier, c *cache.Cache) *Root {
	r := Root{
		Verifier:         verifier,
		RemindDuration:   DefaultRemindDuration,
		InactiveDuration: DefaultInactiveDuration,
		Cache:            c,
		FailStatus:       DefaultLoginFailStatus,
		CookieName:       DefaultCookieName,
	}
	return &r
}
func randMaskedBytes(mask []byte, length int) ([]byte, error) {
	token := make([]byte, length)
	masked := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return masked, err
	}
	l := len(mask)
	for k, v := range token {
		index := int(v) % l
		masked[k] = mask[index]
	}
	return masked, nil
}

type Root struct {
	Verifier         UserVerifier
	RemindDuration   time.Duration
	InactiveDuration time.Duration
	Cache            *cache.Cache
	FailStatus       int
	CookieName       string
	CookiePath       string
}
type ApiResult struct {
	Token string
	Msg   string
}

func (root *Root) Token() (string, error) {
	_, err := root.Cache.GetBytesValue(CacheRefreshedTimeKey)
	if err == cache.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}
	token, err := root.Cache.GetBytesValue(CacheTokenKey)
	if err == cache.ErrNotFound {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return string(token), nil
}
func (root *Root) SetToken(Token string) error {
	t, err := time.Now().MarshalBinary()
	if err != nil {
		return nil
	}
	err = root.Cache.SetBytesValue(CacheRefreshedTimeKey, t, root.InactiveDuration)
	if err != nil {
		return nil
	}
	err = root.Cache.SetBytesValue(CacheTokenKey, []byte(Token), root.RemindDuration)
	return err
}

func (root *Root) Login(r *http.Request) (bool, error) {
	success, err := root.Verifier.Verify(r)
	return success, err
}
func (root *Root) LoginMiddleware(langingUrl string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		success, err := root.Login(r)
		if err != nil {
			panic(err)
		}
		if success {
			err = root.setCookie(w)
			if err != nil {
				panic(err)
			}
			if langingUrl != "" {
				http.Redirect(w, r, langingUrl, http.StatusFound)
			}
		}
	}
}
func (root *Root) setCookie(w http.ResponseWriter) error {
	token, err := root.Token()
	if err != nil {
		return err
	}
	cookie := http.Cookie{
		Name:     root.CookieName,
		Value:    token,
		Path:     root.CookiePath,
		Secure:   false,
		HttpOnly: true,
	}
	if root.RemindDuration >= 0 {
		cookie.Expires = time.Now().Add(root.RemindDuration)
	}
	http.SetCookie(w, &cookie)
	return nil
}
func (root *Root) LoginJSONApi(setCookie bool, msg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		success, err := root.Login(r)
		if err != nil {
			panic(err)
		}
		token, err := root.Regenerate()
		if err != nil {
			panic(err)
		}
		if success {
			result := ApiResult{}
			result.Msg = msg
			result.Token = token
			bytes, err := json.Marshal(result)
			if err != nil {
				panic(err)
			}
			if setCookie {
				err := root.setCookie(w)
				if err != nil {
					panic(err)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)

		} else {
			http.Error(w, "Login fail", root.FailStatus)
		}

	}
}
func (root *Root) Regenerate() (string, error) {
	tokenLock.Lock()
	defer tokenLock.Unlock()
	bytes, err := randMaskedBytes(TokenMask, Tokenlength)
	if err != nil {
		return "", err
	}
	token := string(bytes)
	err = root.SetToken(token)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (root *Root) Verify(str string) (bool, error) {
	token, err := root.Token()
	if err != nil {
		return false, nil
	}
	if token == "" {
		return false, nil
	}
	result := str == token
	return result, nil
}

func (root *Root) VerifyCookieMiddleware(loginUrl string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var success bool
		token, err := r.Cookie(root.CookieName)
		if err == http.ErrNoCookie {
			success = false
		} else if err != nil {
			panic(err)
		} else {
			success, err = root.Verify(token.Value)
			if err != nil {
				panic(err)
			}
		}
		if !success {
			if loginUrl != "" {
				http.Redirect(w, r, loginUrl, http.StatusFound)
			} else {
				http.Error(w, "Forbidden", http.StatusForbidden)
			}
			return
		}
		next(w, r)
	}
}

func (root *Root) InstallRoot(install bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if install == false {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		success, err := root.Verifier.Install(r)
		if err != nil {
			panic(err)
		}
		if success {
			http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

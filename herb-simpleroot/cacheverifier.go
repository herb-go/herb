package simpleroot

import (
	"crypto/sha256"
	"net/http"

	"encoding/json"

	"github.com/herb-go/herb/cache"
)

var DefaultUserKey = "admin"
var DefaulHashLength = 64

type User struct {
	Username string
	Hash     string
	Password string
}

func (u *User) RegenerateHash() error {
	token, err := cache.RandMaskedBytes(cache.TokenMask, DefaulHashLength)
	if err != nil {
		return err
	}
	u.Hash = string(token)
	return nil
}
func (u *User) HashPassword(password string) string {
	sum := sha256.Sum256([]byte(u.Hash + password))
	return string(sum[:])
}
func (u *User) Auth(password string) bool {
	return u.Password == u.HashPassword(password)
}

func (u *User) SetPassword(password string) error {
	err := u.RegenerateHash()
	if err != nil {
		return err
	}
	u.Password = u.HashPassword(password)
	return nil
}

type CacheVerifier struct {
	cache *cache.Cache
	Key   string
}

func (v *CacheVerifier) getUser() (User, error) {
	var user User
	err := v.cache.Get(v.Key, &user)
	return user, err
}
func (v *CacheVerifier) saveUser(user User) error {
	return v.cache.Set(v.Key, user, cache.TTLForever)
}

type UserForm struct {
	Username string
	Password string
}

func NewCacheVerifier(c *cache.Cache) *CacheVerifier {
	v := CacheVerifier{
		cache: c,
		Key:   DefaultUserKey,
	}
	return &v
}
func (v *CacheVerifier) Verify(r *http.Request) (Verified bool, err error) {
	user, err := v.getUser()
	if err != nil {
		return false, err
	}
	var uf UserForm
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&uf)
	if err != nil {
		return false, err
	}
	if uf.Username == user.Username && user.Auth(uf.Password) {
		return true, nil
	}
	return false, nil
}
func (v *CacheVerifier) Install(r *http.Request) (success bool, err error) {
	_, err = v.getUser()
	if err != cache.ErrNotFound {
		return false, err
	}
	var uf UserForm
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&uf)
	if err != nil {
		return false, err
	}
	if uf.Username == "" || uf.Password == "" {
		return false, nil
	}
	var user User
	user.Username = uf.Username
	user.SetPassword(uf.Password)
	err = v.saveUser(user)
	if err != nil {
		return false, err
	}
	return true, nil

}

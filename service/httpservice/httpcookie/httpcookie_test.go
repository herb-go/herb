package httpcookie

import (
	"net/http"
	"testing"
)

func TestConfig(t *testing.T) {
	c := &Config{
		Name:     "testname",
		Path:     "/testpath",
		Domain:   "testdomain",
		Secure:   true,
		HTTPOnly: true,
	}
	cookie := c.CreateCookieWithValue("testvalue")
	if cookie.Name != "testname" || cookie.Path != c.Path || cookie.Domain != c.Domain || cookie.Secure != true || cookie.HttpOnly != true || cookie.Value != "testvalue" {
		t.Fatal(cookie)
	}
	c.SameSite = "wrong"
	cookie = c.CreateCookie()
	if cookie.SameSite != 0 {
		t.Fatal(cookie)
	}
	c.SameSite = SameSiteNameDefault
	cookie = c.CreateCookie()
	if cookie.SameSite != http.SameSiteDefaultMode {
		t.Fatal(cookie)
	}
	c.SameSite = SameSiteNameLax
	cookie = c.CreateCookie()
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatal(cookie)
	}
	c.SameSite = SameSiteNameStrict
	cookie = c.CreateCookie()
	if cookie.SameSite != http.SameSiteStrictMode {
		t.Fatal(cookie)
	}
}

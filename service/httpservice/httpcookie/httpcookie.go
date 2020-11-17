package httpcookie

import "net/http"

//SameSiteName same site name type
type SameSiteName string

//SameSiteNameDefault same site name for default mode
const SameSiteNameDefault = SameSiteName("default")

//SameSiteNameLax same site name for lax mode
const SameSiteNameLax = SameSiteName("lax")

//SameSiteNameStrict same site name for strict mode
const SameSiteNameStrict = SameSiteName("strict")

//SameSiteNameMap same site name map to http.SameSite
var SameSiteNameMap = map[SameSiteName]http.SameSite{
	SameSiteNameDefault: http.SameSiteDefaultMode,
	SameSiteNameLax:     http.SameSiteLaxMode,
	SameSiteNameStrict:  http.SameSiteStrictMode,
}

//Config cookie creator config
type Config struct {
	//Name cookie name
	Name string
	//Path cookie path
	Path string
	//Doain cookie domain
	Domain string
	//Secure cookie secure
	Secure bool
	//HTTPOnly cookie http only
	HTTPOnly bool
	//SameSite cookie same site mdoe
	SameSite SameSiteName
}

//CreateCookie create new cookie .
func (c *Config) CreateCookie() *http.Cookie {
	cookie := &http.Cookie{}
	cookie.Name = c.Name
	cookie.Path = c.Path
	cookie.Domain = c.Domain
	cookie.Secure = c.Secure
	cookie.HttpOnly = c.HTTPOnly
	cookie.SameSite = SameSiteNameMap[c.SameSite]
	return cookie
}

//CreateCookieWithValue create new cookie with given value.
func (c *Config) CreateCookieWithValue(value string) *http.Cookie {
	cookie := c.CreateCookie()
	cookie.Value = value
	return cookie
}

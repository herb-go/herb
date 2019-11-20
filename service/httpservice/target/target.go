package target

import (
	"io"
	"net/http"
)

type Target interface {
	//RequestMethod return request method
	RequestMethod() string
	//RequestURL return request url
	RequestURL() string
	//RequestBody return request body
	RequestBody() io.Reader
	//RequestBuilders return request builders
	RequestBuilders() []func(*http.Request) error
}

type MutableTarget interface {
	Target
	//SetRequestMethod set request method
	SetRequestMethod(string)
	//SetRequesetURL set request url
	SetRequesetURL(string)
	//SetRequestBody set request body
	SetRequestBody(io.Reader)
	//SetRequestBuilders set request builders
	SetRequestBuilders([]func(*http.Request) error)
}

type PlainTarget struct {
	Method   string
	URL      string
	Body     io.Reader
	Builders []func(*http.Request) error
}

func (t *PlainTarget) RequestMethod() string {
	return t.Method
}
func (t *PlainTarget) RequestURL() string {
	return t.URL
}
func (t *PlainTarget) RequestBody() io.Reader {
	return t.Body
}
func (t *PlainTarget) RequestBuilders() []func(*http.Request) error {
	return t.Builders
}
func (t *PlainTarget) SetRequestMethod(v string) {
	t.Method = v
}
func (t *PlainTarget) SetRequesetURL(v string) {
	t.URL = v
}
func (t *PlainTarget) SetRequestBody(v io.Reader) {
	t.Body = v
}
func (t *PlainTarget) SetRequestBuilders(v []func(*http.Request) error) {
	t.Builders = v
}
func (t *PlainTarget) NewRequest() (*http.Request, error) {
	req, err := http.NewRequest(t.Method, t.URL, t.Body)
	if err != nil {
		return nil, err
	}
	for k := range t.Builders {
		err = t.Builders[k](req)
		if err != nil {
			return nil, err
		}
	}
	return req, err
}
func New() *PlainTarget {
	return &PlainTarget{}
}
func Copy(dst MutableTarget, src Target) {
	dst.SetRequesetURL(src.RequestURL())
	dst.SetRequestMethod(src.RequestMethod())
	dst.SetRequestBody(src.RequestBody())
	b := dst.RequestBuilders()
	clone := make([]func(*http.Request) error, len(b))
	copy(clone, b)
}
func Clone(t Target) *PlainTarget {
	pt := New()
	if t != nil {
		Copy(pt, t)
	}
	return pt
}

func BuildPlainTarget(pt *PlainTarget, b ...Builder) error {
	var err error
	for k := range b {
		pt.Method, pt.URL, pt.Body, pt.Builders, err = b[k].BuildTarget(pt.Method, pt.URL, pt.Body, pt.Builders)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRequest(t Target, b ...Builder) (req *http.Request, err error) {
	pt := Clone(t)
	err = BuildPlainTarget(pt)
	if err != nil {
		return nil, err
	}
	return pt.NewRequest()
}

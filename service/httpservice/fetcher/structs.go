package fetcher

import (
	"bytes"
	"io"
	"net/http"
)

type EndPoint struct {
	Parent         TargetGetter
	TargetBuilders []TargetBuilder
}

func NewPlainResult() *PlainResult {
	return &PlainResult{}
}

type PlainResult struct {
	*http.Response
}

//FromResponse convert result From response
func (r *PlainResult) FromResponse(resp *http.Response) error {
	r.Response = resp
	return nil
}

//ToResponse convert result to response
func (r *PlainResult) ToResponse() (*http.Response, error) {
	return r.Response, nil
}
func (r *PlainResult) Error() string {

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

func NewPlainTarget() *PlainTarget {
	return &PlainTarget{}
}

func MakeTargetPlain(t TargetGetter) *PlainTarget {
	pt := NewPlainTarget()
	if t == nil {
		return pt
	}
	pt.Method = t.RequestMethod()
	pt.URL = t.RequestURL()
	pt.Builders = t.RequestBuilders()
	pt.Body = t.RequestBody()
	return pt
}

type Method string

func (m Method) BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error) {
	return string(m), url, body, builders, nil
}

type URL string

func (u URL) BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error) {
	return method, string(u), body, builders, nil
}

type Body []byte

func (b Body) BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error) {
	return method, url, bytes.NewReader(b), builders, nil
}

type MarshalerBody struct {
	reader io.Reader
	err    error
}

func (b *MarshalerBody) Read(p []byte) (n int, err error) {
	if b.err != nil {
		return 0, b.err
	}
	return b.reader.Read(p)
}
func (b *MarshalerBody) BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error) {
	return method, url, b, builders, nil
}
func NewMarshalerBody(m func(v interface{}) ([]byte, error), v interface{}) *MarshalerBody {
	b := &MarshalerBody{}
	bs, err := m(v)
	if err != nil {
		b.err = err
	} else {
		b.reader = bytes.NewReader(bs)
	}
	return b
}

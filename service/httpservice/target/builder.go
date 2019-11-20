package target

import (
	"bytes"
	"io"
	"net/http"
)

type Builder interface {
	//BuildTarget build given request method ,request url request body ,request builders and return any error raised
	BuildTarget(
		method string,
		url string,
		body io.Reader,
		builders []func(*http.Request) error,
	) (
		newmethod string,
		newurl string,
		newbody io.Reader,
		newbuilders []func(*http.Request) error,
		err error,
	)
}
type Builders []Builder

func (b *Builders) BuildTarget(
	method string,
	url string,
	body io.Reader,
	builders []func(*http.Request) error,
) (
	newmethod string,
	newurl string,
	newbody io.Reader,
	newbuilders []func(*http.Request) error,
	err error,
) {
	newmethod = method
	newurl = url
	newbody = body
	newbuilders = builders
	for k := range *b {
		newmethod, newurl, newbody, newbuilders, err = (*b)[k].BuildTarget(newmethod, newurl, newbody, newbuilders)
		if err != nil {
			return "", "", nil, nil, err
		}
	}
	return newmethod, newurl, newbody, newbuilders, nil
}

func (b *Builders) AppendBuilder(builders ...Builder) {
	*b = append(*b, builders...)
}
func NewBuilders(b ...Builder) *Builders {
	builders := Builders(b)
	return &builders
}

type Method string

//BuildTarget build given request method ,request url request body ,request builders and return any error raised
func (m Method) BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error) {
	return string(m), url, body, builders, nil
}

type URL string

//BuildTarget build given request method ,request url request body ,request builders and return any error raised
func (u URL) BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error) {
	return method, string(u), body, builders, nil
}

type Body []byte

//BuildTarget build given request method ,request url request body ,request builders and return any error raised
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

//BuildTarget build given request method ,request url request body ,request builders and return any error raised
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

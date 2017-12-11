package fetch

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var ErrMsgLengthLimit = 512

var DefaultTimeout = 120
var DefaultMaxIdleConns = 100
var DefaultIdleConnTimeout = 120 * time.Second
var DefaultTLSHandshakeTimeout = 30 * time.Second

type Service struct {
	TimeoutInSecond     int
	MaxIdleConns        int
	IdleConnTimeout     int
	TLSHandshakeTimeout int
	ProxyURL            string
	proxyCache          map[string]func(*http.Request) (*url.URL, error)
}

var DefaultService = &Service{}

func (s *Service) getProxy(index string) func(*http.Request) (*url.URL, error) {
	if index == "" {
		return http.ProxyFromEnvironment
	}
	p, ok := s.proxyCache[index]
	if ok == false {
		url, err := url.Parse(index)
		if err != nil {
			p = nil
		} else {
			p = http.ProxyURL(url)
		}
	}
	return p
}
func (s *Service) Client() *Client {
	var cs *Service
	if s != nil {
		cs = s
	} else {
		cs = DefaultService
	}

	timeout := cs.TimeoutInSecond
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	c := http.Client{
		Transport: s.getTransport(),
		Timeout:   time.Duration(cs.TimeoutInSecond) * time.Second,
	}
	return &Client{Client: &c}
}
func (s *Service) getTransport() *http.Transport {
	var maxIdleCoons = s.MaxIdleConns
	if maxIdleCoons == 0 {
		maxIdleCoons = DefaultMaxIdleConns
	}
	var idleConnTimeout = time.Duration(s.IdleConnTimeout) * time.Second
	if idleConnTimeout == 0 {
		idleConnTimeout = DefaultIdleConnTimeout
	}
	var tlsHandshakeTimeout = time.Duration(s.TLSHandshakeTimeout) * time.Second
	if tlsHandshakeTimeout == 0 {
		tlsHandshakeTimeout = DefaultTLSHandshakeTimeout
	}
	return &http.Transport{
		Proxy:               s.getProxy(s.ProxyURL),
		MaxIdleConns:        maxIdleCoons,
		IdleConnTimeout:     idleConnTimeout,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
}
func (s *Service) Fetch(req *http.Request) (*Result, error) {
	return s.Client().Fetch(req)
}

type Client struct {
	*http.Client
}

func (c *Client) Fetch(req *http.Request) (*Result, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := Result{
		Response:    resp,
		BodyContent: bodyContent,
	}
	return &result, nil
}

type Result struct {
	*http.Response
	BodyContent []byte
}

func (r *Result) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.BodyContent, &v)
}
func (r *Result) UnmarshalXML(v interface{}) error {
	return xml.Unmarshal(r.BodyContent, &v)
}
func (r Result) Error() string {
	msg := fmt.Sprintf("http error [ %s ] %s : %s", r.Response.Request.URL.String(), r.Status, string(r.BodyContent))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

func (r *Result) NewAPICodeErr(code interface{}) *APICodeErr {
	return NewAPICodeErr(r.Response.Request.URL.String(), code, r.BodyContent)

}
func GetErrorStatusCode(err error) int {
	r, ok := err.(Result)
	if ok {
		return r.StatusCode
	}
	return 0
}

func NewAPICodeErr(url string, code interface{}, content []byte) *APICodeErr {
	return &APICodeErr{
		URI:     url,
		Code:    fmt.Sprint(code),
		Content: content,
	}
}

type APICodeErr struct {
	URI     string
	Code    string
	Content []byte
}

func (r APICodeErr) Error() string {
	msg := fmt.Sprintf("api error [ %s] code %d : %s", r.URI, r.Code, string(r.Content))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

func GetAPIErrCode(err error) string {
	r, ok := err.(APICodeErr)
	if ok {
		return r.Code
	}
	return ""

}

func CompareApiErrCode(err error, code interface{}) bool {
	return GetAPIErrCode(err) == fmt.Sprint(code)
}

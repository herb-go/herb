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

//ErrMsgLengthLimit max error message length
var ErrMsgLengthLimit = 512

//DefaultTimeout default client timeout
var DefaultTimeout = 120

//DefaultMaxIdleConns default client max idle conns
var DefaultMaxIdleConns = 20

//DefaultIdleConnTimeout default client idle conn timeout.
var DefaultIdleConnTimeout = 120 * time.Second

//DefaultTLSHandshakeTimeout default client tls handshake time out.
var DefaultTLSHandshakeTimeout = 30 * time.Second

//Clients fetch clients struct.
type Clients struct {
	//TimeoutInSecond timeout in secound
	//Default value is 120.
	TimeoutInSecond int
	//MaxIdleConns max idel conns.
	//Default value is 20
	MaxIdleConns int
	//IdleConnTimeoutInSecond idel conn timeout in second.
	//Default value is 120
	IdleConnTimeoutInSecond int
	//TLSHandshakeTimeoutInSecond tls handshake timeout in secound.
	//default value is 30.
	TLSHandshakeTimeoutInSecond int
	//ProxyURL proxy url.
	//If set to empty string,clients will not use proxy.
	//Default value is empty string.
	ProxyURL   string
	proxyCache map[string]func(*http.Request) (*url.URL, error)
}

//DefaultClients default fetch clients.
var DefaultClients = &Clients{}

func (clients *Clients) getProxy(index string) func(*http.Request) (*url.URL, error) {
	if index == "" {
		return http.ProxyFromEnvironment
	}
	p, ok := clients.proxyCache[index]
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

//Client get a client from fetch clients.
//This client has all method than http.Client has.
//This client can use fetch method to get a fetch result.
func (clients *Clients) Client() *Client {
	var cs *Clients
	if clients != nil {
		cs = clients
	} else {
		cs = DefaultClients
	}

	timeout := cs.TimeoutInSecond
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	c := http.Client{
		Transport: cs.getTransport(),
		Timeout:   time.Duration(cs.TimeoutInSecond) * time.Second,
	}
	return &Client{Client: &c}
}
func (clients *Clients) getTransport() *http.Transport {
	var maxIdleCoons = clients.MaxIdleConns
	if maxIdleCoons == 0 {
		maxIdleCoons = DefaultMaxIdleConns
	}
	var idleConnTimeout time.Duration
	if idleConnTimeout == 0 {
		idleConnTimeout = DefaultIdleConnTimeout
	} else {
		idleConnTimeout = time.Duration(clients.IdleConnTimeoutInSecond) * time.Second
	}
	var tlsHandshakeTimeout time.Duration
	if tlsHandshakeTimeout == 0 {
		tlsHandshakeTimeout = DefaultTLSHandshakeTimeout
	} else {
		tlsHandshakeTimeout = time.Duration(clients.TLSHandshakeTimeoutInSecond) * time.Second
	}
	return &http.Transport{
		Proxy:               clients.getProxy(clients.ProxyURL),
		MaxIdleConns:        maxIdleCoons,
		IdleConnTimeout:     idleConnTimeout,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
}

//Fetch fetch a fetch result.
func (clients *Clients) Fetch(req *http.Request) (*Result, error) {
	return clients.Client().Fetch(req)
}

//Client Fetch client
type Client struct {
	*http.Client
}

//Fetch fetch a fetch result.
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

//Result fetch result.
//Result will read all bytes form body,store data in BodyContent field,and close the body.
//Result can be used as error directly.
type Result struct {
	*http.Response
	BodyContent []byte
}

//UnmarshalAsJSON unmarshal body content as JSON.
func (r *Result) UnmarshalAsJSON(v interface{}) error {
	return json.Unmarshal(r.BodyContent, v)
}

//UnmarshalAsXML unmarshal body content as XML.
func (r *Result) UnmarshalAsXML(v interface{}) error {
	return xml.Unmarshal(r.BodyContent, v)
}

//Error result can used as a error which return request url,request status,request content.
//Error max length is ErrMsgLengthLimit.
func (r Result) Error() string {
	msg := fmt.Sprintf("http error [ %s ] %s : %s", r.Response.Request.URL.String(), r.Status, string(r.BodyContent))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

//NewAPICodeErr make a api code error  which contains a error code.
func (r *Result) NewAPICodeErr(code interface{}) *APICodeErr {
	return NewAPICodeErr(r.Response.Request.URL.String(), code, r.BodyContent)

}

//GetErrorStatusCode get status code form response error.
//Return 0 if error is not a fetch error.
func GetErrorStatusCode(err error) int {
	r, ok := err.(*Result)
	if ok {
		return r.StatusCode
	}
	r2, ok := err.(Result)
	if ok {
		return r2.StatusCode
	}
	return 0
}

//NewAPICodeErr create a new api code error with given url,code,and content.
func NewAPICodeErr(url string, code interface{}, content []byte) *APICodeErr {
	return &APICodeErr{
		URI:     url,
		Code:    fmt.Sprint(code),
		Content: content,
	}
}

//APICodeErr api code error struct.
type APICodeErr struct {
	//URI api uri.
	URI string
	//Code api error code.
	Code string
	//Content api response.
	Content []byte
}

//Error used as a error which return request url,request status,erro code,request content.
//Error max length is ErrMsgLengthLimit.
func (r APICodeErr) Error() string {
	msg := fmt.Sprintf("api error [ %s] code %s : %s", r.URI, r.Code, string(r.Content))
	if len(msg) > ErrMsgLengthLimit {
		msg = msg[:ErrMsgLengthLimit]
	}
	return msg
}

//GetAPIErrCode get api error code form error.
//Return empty string if err is not an ApiCodeErr
func GetAPIErrCode(err error) string {
	r, ok := err.(APICodeErr)
	if ok {
		return r.Code
	}
	r2, ok := err.(*APICodeErr)
	if ok {
		return r2.Code
	}
	return ""

}

//CompareAPIErrCode if check error is an ApiCodeErr with given api err code.
func CompareAPIErrCode(err error, code interface{}) bool {
	return GetAPIErrCode(err) == fmt.Sprint(code)
}

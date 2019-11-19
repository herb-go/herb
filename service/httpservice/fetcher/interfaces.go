package fetcher

import (
	"io"
	"net/http"
)

type RequestBuilder func(*http.Request) error

type TargetGetter interface {
	//RequestMethod return request method
	RequestMethod() string
	//RequestURL return request url
	RequestURL() string
	//RequestBody return request body
	RequestBody() io.Reader
	//RequestBuilders return request builders
	RequestBuilders() []func(*http.Request) error
}

type TargetSetter interface {
	//SetRequestMethod set request method
	SetRequestMethod(string)
	//SetRequesetURL set request url
	SetRequesetURL(string)
	//SetRequestBody set request body
	SetRequestBody(io.Reader)
	//SetRequestBuilders set request builders
	SetRequestBuilders([]func(*http.Request) error)
}

type Target interface {
	TargetGetter
	TargetSetter
}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type Result interface {
	//FromResponse convert result From response
	FromResponse(*http.Response) error
	//ToResponse convert result to response
	ToResponse() (*http.Response, error)
	Error() string
}

type TargetBuilder interface {
	//BuildRequest builde given request and return any error raised
	BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error)
}

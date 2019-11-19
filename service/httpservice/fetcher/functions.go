package fetcher

import (
	"io"
	"net/http"
)

func NewRequest(t TargetGetter, b ...TargetBuilder) (req *http.Request, err error) {
	var method string
	var url string
	var body io.Reader
	var builders []func(*http.Request) error

	for k := range b {
		method, url, body, builders, err = b[k].BuildTarget(method, url, body, builders)
		if err != nil {
			return
		}
	}
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k := range builders {
		err = builders[k](req)
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func Fetch(result Result, client Client, t TargetGetter, b ...TargetBuilder) error {
	if client == nil {
		client = http.DefaultClient
	}
	req, err := NewRequest(t, b...)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return result.FromResponse(resp)
}

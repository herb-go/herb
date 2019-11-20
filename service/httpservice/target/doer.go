package target

import "net/http"

var DefaultDoer = http.DefaultClient

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func Do(doer Doer, t Target, b ...Builder) (*http.Response, error) {
	if doer == nil {
		doer = DefaultDoer
	}
	req, err := NewRequest(t, b...)
	if err != nil {
		return nil, err
	}
	return doer.Do(req)
}

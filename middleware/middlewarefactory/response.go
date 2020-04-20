package middlewarefactory

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

//ResponseMiddlewareConfig response middleware config struct
type ResponseMiddlewareConfig struct {
	StatusCode int
	Header     http.Header
	Body       *string
}

//ResponseMiddleware response middleware config
type ResponseMiddleware struct {
	StatusCode int
	Header     http.Header
	Body       *string
}

//ServeMiddleware serve as middleware
func (m *ResponseMiddleware) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	for field := range m.Header {
		for key := range m.Header[field] {
			w.Header().Add(field, m.Header[field][key])
		}
	}
	if m.StatusCode != 0 {
		w.WriteHeader(m.StatusCode)
	}
	var body string
	if m.Body == nil {
		body = http.StatusText(m.StatusCode)

	} else {
		body = *m.Body
	}
	_, err := w.Write([]byte(body))
	if err != nil {
		panic(err)
	}
	return
}

//NewResponseFactory create new response factory
var NewResponseFactory = func() Factory {
	return FactoryFunc(func(name string, loader func(v interface{}) error) (middleware.Middleware, error) {
		m := &ResponseMiddleware{}
		err := loader(m)
		if err != nil {
			return nil, err
		}
		return m.ServeMiddleware, nil
	})
}

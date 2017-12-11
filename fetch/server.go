package fetch

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
)

type Server struct {
	Host    string
	Headers http.Header
}

func (s *Server) EndPoint(method string, path string) *EndPoint {
	return &EndPoint{
		Server: s,
		Method: method,
		Path:   path,
	}
}
func (s *Server) NewRequest(method string, path string, params url.Values, body []byte) (*http.Request, error) {
	u, err := url.Parse(s.Host + path)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	if params != nil {
		for k, vs := range params {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if s.Headers != nil {
		for k, vs := range s.Headers {
			for _, v := range vs {
				req.Header.Add(k, v)
			}
		}
	}
	return req, nil
}

func (s *Server) NewJSONRequest(method string, path string, params url.Values, v interface{}) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return s.NewRequest(method, path, params, b)
}

func (s *Server) NewXMLRequest(method string, path string, params url.Values, v interface{}) (*http.Request, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return s.NewRequest(method, path, params, b)
}

type EndPoint struct {
	Server *Server
	Path   string
	Method string
}

func (e *EndPoint) NewRequest(params url.Values, body []byte) (*http.Request, error) {
	return e.Server.NewRequest(e.Method, e.Path, params, body)
}

func (e *EndPoint) NewJSONRequest(params url.Values, v interface{}) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return e.NewRequest(params, b)
}

func (e *EndPoint) NewXMLRequest(params url.Values, v interface{}) (*http.Request, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return e.NewRequest(params, b)
}

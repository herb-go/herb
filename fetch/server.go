package fetch

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
)

//Server api server struct
type Server struct {
	//Host server base host
	Host string
	//Headers headers which will be sent every request.
	Headers http.Header
}

//EndPoint create a new api server endpoint with given method and path
func (s *Server) EndPoint(method string, path string) *EndPoint {
	return &EndPoint{
		Server: s,
		Method: method,
		Path:   path,
	}
}

//NewRequest create a new http.request with given method,path,params,and body.
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

//NewJSONRequest create a new http.request with given method,path,params,and body encode by JSON.
func (s *Server) NewJSONRequest(method string, path string, params url.Values, v interface{}) (*http.Request, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return s.NewRequest(method, path, params, b)
}

//NewXMLRequest create a new http.request with given method,path,params,and body encode by XML.
func (s *Server) NewXMLRequest(method string, path string, params url.Values, v interface{}) (*http.Request, error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	return s.NewRequest(method, path, params, b)
}

//EndPoint api server endpoint struct
//Endpoint should be created by api server's EndPoint method
type EndPoint struct {
	Server *Server
	Path   string
	Method string
}

//NewRequest create a new http.request to end point with given params,and body.
func (e *EndPoint) NewRequest(params url.Values, body []byte) (*http.Request, error) {
	return e.Server.NewRequest(e.Method, e.Path, params, body)
}

//NewJSONRequest create a new http.request to end point with given params,and body encode by JSON.
func (e *EndPoint) NewJSONRequest(params url.Values, v interface{}) (*http.Request, error) {
	return e.Server.NewJSONRequest(e.Method, e.Path, params, v)
}

//NewXMLRequest create a new http.request to end point with given params,and body encode by XML.
func (e *EndPoint) NewXMLRequest(params url.Values, v interface{}) (*http.Request, error) {
	return e.Server.NewXMLRequest(e.Method, e.Path, params, v)
}

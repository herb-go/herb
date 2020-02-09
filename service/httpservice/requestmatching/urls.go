package requestmatching

import (
	"net/http"
	"strings"
)

//PathData path data type
type PathData map[string]bool

//Paths paths pattern struct
type Paths map[string]PathData

//Add add url to paths pattern
//Any string before first "/" will be used as hostname
func (p Paths) Add(url string) {
	host, path := splitURL(url)
	data := p[host]
	if data == nil {
		data = PathData{}
		p[host] = data
	}
	p[host][path] = true
}

//MatchRequest match request.
//Return result and any error if raised.
func (p Paths) MatchRequest(r *http.Request) (bool, error) {
	var data PathData
	if len(p) == 0 {
		return true, nil
	}
	urlpath := strings.ToLower(r.URL.Path)
	h := r.Host
	if h != "" {
		data = p[h]
		if data != nil {
			if data[urlpath] == true {
				return true, nil
			}
		}
	}
	data = p[""]
	if data == nil {
		return false, nil
	}

	return data[r.URL.Path], nil
}

//NewPaths create new paths pattern.
func NewPaths() *Paths {
	return &Paths{}
}
func splitURL(url string) (host string, path string) {
	url = strings.ToLower(url)
	if url[0] == '/' {
		path = url
	} else {
		splited := strings.SplitN(url, "/", 2)
		host = splited[0]
		if len(splited) == 2 {
			path = "/" + splited[1]
		} else {
			path = "/"
		}
	}
	return host, path
}

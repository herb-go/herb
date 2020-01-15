package requestmatching

import (
	"net/http"
	"path/filepath"
	"strings"
)

//Exts request exts pattern
type Exts map[string]bool

//MatchRequest match request.
//Request file ext will be converted to lower when matching
//Return result and any error if raised
func (e *Exts) MatchRequest(r *http.Request) (bool, error) {
	if len(*e) == 0 {
		return true, nil
	}
	return (*e)[strings.ToLower(filepath.Ext(r.URL.Path))], nil
}

//Add add ext to pattern
//File ext will be be converted to lower.
func (e *Exts) Add(ext string) {
	(*e)[strings.ToLower(ext)] = true
}

//NewExts create new exts pattern
func NewExts() *Exts {
	return &Exts{}
}

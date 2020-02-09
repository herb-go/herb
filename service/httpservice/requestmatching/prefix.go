package requestmatching

import (
	"net/http"
	"strings"
)

//PrefixData prefix data struct.
type PrefixData []string

//Has check is requesturi is start with any record in prefix data.
func (p PrefixData) Has(requesturi string) bool {
	for k := range p {
		if strings.HasPrefix(requesturi, p[k]) {
			return true
		}
	}
	return false
}

//Prefixs prefixs pattern struct
type Prefixs map[string]PrefixData

//MatchRequest match request.
//Return result and any error if raised.
func (p Prefixs) MatchRequest(r *http.Request) (bool, error) {
	var data PrefixData
	if len(p) == 0 {
		return true, nil
	}
	urlpath := strings.ToLower(r.URL.Path)
	h := r.Host
	if h != "" {
		data = p[h]
		if data != nil && data.Has(urlpath) {
			return true, nil

		}
	}
	data = p[""]
	if data == nil {
		return false, nil
	}

	return data.Has(r.URL.Path), nil
}

//Add add url to prefixs
//Any string before first "/" will be used as hostname
func (p Prefixs) Add(url string) {
	host, prefix := splitURL(url)
	data := p[host]
	if data == nil {
		data = PrefixData{}
		p[host] = data
	}
	p[host] = append(p[host], prefix)
}

//NewPrefixs create new prefixs pattern.
func NewPrefixs() *Prefixs {
	return &Prefixs{}
}

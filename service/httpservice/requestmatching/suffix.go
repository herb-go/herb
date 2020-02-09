package requestmatching

import (
	"net/http"
	"strings"
)

//SuffixData prefix data struct.
type SuffixData []string

//Has check is requesturi is start with any record in prefix data.
func (p SuffixData) Has(requesturi string) bool {
	for k := range p {
		if strings.HasSuffix(requesturi, p[k]) {
			return true
		}
	}
	return false
}

//Suffixs prefixs pattern struct
type Suffixs map[string]SuffixData

//MatchRequest match request.
//Return result and any error if raised.
func (p Suffixs) MatchRequest(r *http.Request) (bool, error) {
	var data SuffixData
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
func (p Suffixs) Add(url string) {
	host, suffix := splitURL(url)
	//remove first "/"
	suffix = suffix[1:]
	data := p[host]
	if data == nil {
		data = SuffixData{}
		p[host] = data
	}
	p[host] = append(p[host], suffix)
}

//NewSuffixs create new prefixs pattern.
func NewSuffixs() *Suffixs {
	return &Suffixs{}
}

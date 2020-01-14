package requestmatching

import (
	"net/http"
	"strings"
)

type Paths map[string]bool

type PathsWithHost map[string]Paths

func (p PathsWithHost) MatchRequest(r *http.Request) (bool, error) {
	var paths Paths
	if len(p) == 0 {
		return true, nil
	}
	h := r.Host
	if h != "" {
		paths = p[h]
		if paths != nil {
			if paths[r.RequestURI] == true {
				return true, nil
			}
		}
	}
	paths = p[""]
	if paths == nil {
		return false, nil
	}

	return paths[r.RequestURI], nil
}

type Prefixs []string

func (p Prefixs) has(requesturi string) bool {
	for k := range p {
		if strings.HasPrefix(requesturi, p[k]) {
			return true
		}
	}
	return false
}

type PrefixsWithHost map[string]Prefixs

func (p PrefixsWithHost) MatchRequest(r *http.Request) (bool, error) {
	var prefixs Prefixs
	if len(p) == 0 {
		return true, nil
	}
	h := r.Host
	if h != "" {
		prefixs = p[h]
		if prefixs != nil && prefixs.has(r.RequestURI) {
			return true, nil

		}
	}
	prefixs = p[""]
	if prefixs == nil {
		return false, nil
	}

	return prefixs.has(r.RequestURI), nil
}

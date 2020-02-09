package requestmatching

import (
	"net/http"
	"strings"
)

type Keywords []string

func (k *Keywords) Add(keyword string) {
	*k = append(*k, strings.ToLower(keyword))
}

func (k *Keywords) MatchRequest(r *http.Request) (bool, error) {
	if len(*k) == 0 {
		return true, nil
	}
	uri := strings.ToLower(r.URL.RequestURI())
	for i := range *k {
		if strings.Index(uri, (*k)[i]) > -1 {
			return true, nil
		}
	}
	return false, nil
}

func NewKeywords() *Keywords {
	return &Keywords{}
}

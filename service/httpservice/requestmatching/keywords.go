package requestmatching

import (
	"net/http"
	"strings"
)

//Keywords keywords pattern
type Keywords []string

//Add add keyword to keywords.
//Method will be converted to lower.
func (k *Keywords) Add(keyword string) {
	*k = append(*k, strings.ToLower(keyword))
}

//MatchRequest match request.
//Return result and any error if raised.
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

//NewKeywords create new keywords
func NewKeywords() *Keywords {
	return &Keywords{}
}

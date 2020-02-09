package requestmatching

import (
	"net/http"
	"regexp"
)

//RegExps regexps pattern
type RegExps []*regexp.Regexp

//MatchRequest match request.
//Return result and any error if raised.
func (r *RegExps) MatchRequest(req *http.Request) (bool, error) {
	if len(*r) == 0 {
		return true, nil
	}
	for i := range *r {
		if (*r)[i].Match([]byte(req.URL.Path)) {
			return true, nil
		}
	}
	return false, nil
}

//Add add pattern to regexps.
//Return any error if raised.
func (r *RegExps) Add(pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	*r = append(*r, re)
	return nil
}

//NewRegExps create new regexps pattern
func NewRegExps() *RegExps {
	return &RegExps{}
}

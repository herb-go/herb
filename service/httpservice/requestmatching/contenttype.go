package requestmatching

import (
	"net/http"
	"strings"
)

//ContentTypes conttenttypes pattern
type ContentTypes map[string]bool

//MatchRequest match request.
//Return result and any error if raised
func (t *ContentTypes) MatchRequest(r *http.Request) (bool, error) {
	if len(*t) == 0 {
		return true, nil
	}
	contenttype := strings.TrimSpace(strings.SplitN(r.Header.Get("Content-Type"), ";", 2)[0])
	return (*t)[contenttype], nil
}

//Add add content type to pattern
func (t *ContentTypes) Add(contenttype string) {
	(*t)[contenttype] = true
}

//NewContentTypes create new content types pattern
func NewContentTypes() *ContentTypes {
	return &ContentTypes{}
}

package requestmatching

import (
	"fmt"
	"net/http"
	"strings"
)

//HeaderData header data stuct
type HeaderData struct {
	Key   string
	Value string
}

//Headers header pattern
type Headers []HeaderData

//Add add header to pattern.
//Return error if header is not validated("Field: value" form)
func (h *Headers) Add(header string) error {
	v := strings.SplitN(header, ":", 2)
	if len(v) < 2 {
		return fmt.Errorf("%w : \"%s\"", ErrHeaderNotValidated, header)
	}
	(*h) = append(*h, HeaderData{
		Key:   strings.TrimSpace(v[0]),
		Value: strings.TrimSpace(v[1]),
	})
	return nil
}

//MatchRequest match request.
//Return result and any error if raised.
func (h *Headers) MatchRequest(r *http.Request) (bool, error) {
	if len(*h) == 0 {
		return true, nil
	}
	for k := range *h {
		if r.Header.Get((*h)[k].Key) == (*h)[k].Value {
			return true, nil
		}
	}
	return false, nil
}

//NewHeaders create new headers pattern.
func NewHeaders() *Headers {
	return &Headers{}
}

package httpinfo

import (
	"net/http"
	"testing"
)

var validatorNilRequest = ValidatorFunc(func(r *http.Request, resp *Response) (bool, error) {
	return r == nil, nil
})

func TestBuffer(t *testing.T) {

}

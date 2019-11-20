package target

import (
	"net/http"
)

var TextErrorContentLengthLimit = 256

type TextResponseError struct {
	*http.Response
}

func (e *TextResponseError) Error() string {
	return ""
}

func NewError(resp *http.Response) *TextResponseError {
	return &TextResponseError{
		resp,
	}
}

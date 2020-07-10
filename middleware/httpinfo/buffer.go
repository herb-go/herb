package httpinfo

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseBuffer struct {
	request   *http.Request
	Error     error
	validator Validator
	discarded bool
	writer    io.Writer
	buffer    *bytes.Buffer
}

func (b *ResponseBuffer) Write([]byte) {
	if b.discarded != false {
		return
	}
}
func (b *ResponseBuffer) Discard() {
	b.discarded = true
}
func (b *ResponseBuffer) Check(resp *Response) {
	if b.discarded != false {
		return
	}
	ok, err := b.validator.Validate(b.request, resp)
	if err != nil {
		b.Error = err
		b.Discard()
		return
	}
	if !ok {
		b.Discard()
	}
}

func NewResponseBuffer() *ResponseBuffer {
	buf := bytes.NewBuffer(nil)
	return &ResponseBuffer{
		validator: ValidatorAlways,
		writer:    buf,
		buffer:    buf,
	}
}

package httpinfo

import (
	"bytes"
	"io"
	"net/http"
)

type Buffer struct {
	request   *http.Request
	Error     error
	checker   Validator
	discarded bool
	writer    io.Writer
	buffer    *bytes.Buffer
}

func (b *Buffer) Write(data []byte) {
	if b.discarded != false {
		return
	}
	_, err := b.writer.Write(data)
	if err != nil {
		b.Error = err
		b.Discard()
	}
}
func (b *Buffer) Discard() {
	b.discarded = true
}
func (b *Buffer) Check(resp *Response) {
	if b.discarded != false {
		return
	}
	ok, err := b.checker.Validate(b.request, resp)
	if err != nil {
		b.Error = err
		b.Discard()
		return
	}
	if !ok {
		b.Discard()
	}
}

func NewBuffer() *Buffer {
	buf := bytes.NewBuffer(nil)
	return &Buffer{
		checker: ValidatorAlways,
		writer:  buf,
		buffer:  buf,
	}
}

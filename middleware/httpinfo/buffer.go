package httpinfo

import (
	"io"
	"net/http"
)

type BufferPipe struct {
	request   *http.Request
	response  *Response
	lasterror error
	checker   Validator
	discarded bool
	writer    io.Writer
}

func (p *BufferPipe) Error() error {
	return p.lasterror
}
func (p *BufferPipe) Write(data []byte) {
	if p.discarded != false {
		return
	}
	_, err := p.writer.Write(data)
	if err != nil {
		p.lasterror = err
		p.Discard()
	}
}
func (p *BufferPipe) Discard() {
	p.discarded = true
}
func (p *BufferPipe) Discarded() bool {
	return p.discarded
}
func (p *BufferPipe) Check() {
	if p.discarded != false {
		return
	}
	ok, err := p.checker.Validate(p.request, p.response)
	if err != nil {
		p.lasterror = err
		p.Discard()
		return
	}
	if !ok {
		p.Discard()
	}
}
func (p *BufferPipe) WithChecker(v Validator) *BufferPipe {
	p.checker = v
	return p
}
func (p *BufferPipe) WithWriter(w io.Writer) *BufferPipe {
	p.writer = w
	return p
}
func NewBufferPipe(req *http.Request, resp *Response) *BufferPipe {
	return &BufferPipe{
		checker:  ValidatorAlways,
		writer:   resp.Buffer,
		request:  req,
		response: resp,
	}
}

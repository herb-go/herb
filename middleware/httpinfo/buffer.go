package httpinfo

import (
	"bytes"
	"io"
	"net/http"
)

type BufferController struct {
	request   *http.Request
	response  *Response
	lasterror error
	checker   Validator
	buffer    *bytes.Buffer
	discarded bool
	writer    io.Writer
	NopController
}

func (c *BufferController) Error() error {
	return c.lasterror
}
func (c *BufferController) Write(data []byte) {
	if c.discarded != false {
		return
	}
	_, err := c.writer.Write(data)
	if err != nil {
		c.lasterror = err
		c.Discard()
	}
}
func (c *BufferController) Discard() {
	c.discarded = true
}
func (c *BufferController) Discarded() bool {
	return c.discarded
}
func (c *BufferController) BeforeWriteHeader() {
	c.check()
}

func (c *BufferController) BeforeWrite() {
	c.check()
}

func (c *BufferController) check() {
	if c.discarded != false {
		return
	}
	ok, err := c.checker.Validate(c.request, c.response)
	if err != nil {
		c.lasterror = err
		c.Discard()
		return
	}
	if !ok {
		c.Discard()
	}
}
func (c *BufferController) WithChecker(v Validator) *BufferController {
	c.checker = v
	return c
}
func (c *BufferController) WithWriter(w io.Writer) *BufferController {
	c.writer = w
	return c
}
func NewBufferController(req *http.Request, resp *Response) *BufferController {
	buf := bytes.NewBuffer(nil)
	return &BufferController{
		checker:  ValidatorAlways,
		writer:   buf,
		buffer:   buf,
		request:  req,
		response: resp,
	}
}

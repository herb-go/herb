package httpinfo

import (
	"bytes"
	"net/http"

	"github.com/herb-go/herb/middleware"
)

//Response standard http response infomation
type Response struct {
	StatusCode    int //StatusCode response status code.Default value 200
	ContentLength int //ContentLength response content length.
	header        http.Header
	Written       bool //Content written
	writer        http.ResponseWriter
	autocommit    bool
	buffer        *bytes.Buffer
	controller    Controller
	locked        bool
}

//NewResponse create new response
func NewResponse() *Response {
	return &Response{
		StatusCode: 200,
		controller: DefaultController,
		autocommit: true,
		header:     http.Header{},
		buffer:     bytes.NewBuffer(nil),
	}
}
func (resp *Response) flushHeader() {
	for field := range resp.header {
		for k := range resp.header[field] {
			resp.writer.Header().Add(field, resp.header[field][k])
		}
	}
}
func (resp *Response) writeHeaderFunc(statusCode int) {
	resp.locked = true
	resp.StatusCode = statusCode
	resp.controller.BeforeWriteHeader()

}

func (resp *Response) writeFunc(data []byte) (int, error) {
	var err error
	var length int
	resp.locked = true
	if !resp.Written {
		resp.Written = true
		if resp.autocommit {
			resp.flushHeader()
			resp.writer.WriteHeader(resp.StatusCode)
		}
	}
	resp.controller.BeforeWrite()
	if resp.autocommit {
		length, err = resp.writer.Write(data)

	} else {
		length, err = resp.buffer.Write(data)
	}
	if err != nil {
		return 0, err
	}
	resp.ContentLength = resp.ContentLength + length
	resp.controller.Write(data)
	return length, nil
}

//Header http response writer header
func (resp *Response) Header() http.Header {
	resp.locked = true
	return resp.header
}

//WrapWriter wrap http response writer
func (resp *Response) WrapWriter(rw http.ResponseWriter) middleware.ResponseWriter {
	resp.writer = rw
	w := middleware.WrapResponseWriter(rw)
	f := w.Functions()
	f.HeaderFunc = resp.Header
	f.WriteHeaderFunc = resp.writeHeaderFunc
	f.WriteFunc = resp.writeFunc
	return w
}

func (resp *Response) UncommittedData() []byte {
	return resp.buffer.Bytes()
}
func (resp *Response) SetUncommittedData(data []byte) {
	resp.buffer = bytes.NewBuffer(data)
}
func (resp *Response) Commit() error {
	if resp.autocommit {
		return nil
	}
	resp.autocommit = true
	resp.flushHeader()
	if !resp.Written {
		return nil
	}
	resp.writer.WriteHeader(resp.StatusCode)
	_, err := resp.writer.Write(resp.buffer.Bytes())
	return err
}

func (resp *Response) Autocommit() bool {
	return resp.autocommit
}

func (resp *Response) Locked() bool {
	return resp.locked
}
func (resp *Response) UpdateController(c Controller) bool {
	if resp.Locked() {
		return false
	}
	resp.controller = c
	return true
}

func (resp *Response) UpdateAutocommit(autocommit bool) bool {
	if resp.Locked() {
		return false
	}
	resp.autocommit = autocommit
	return true
}

func (resp *Response) LastError() error {
	return resp.controller.Error()
}

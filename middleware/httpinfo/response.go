package httpinfo

import (
	"io"
	"net/http"

	"github.com/herb-go/herb/middleware"
)

//Response standard http response infomation
type Response struct {
	StatusCode    int  //StatusCode response status code.Default value 200
	ContentLength int  //ContentLength response content length.
	Written       bool //Content written
	writer        http.ResponseWriter
	buffer        *ResponseBuffer
}

//NewResponse create new response
func NewResponse() *Response {
	return &Response{
		StatusCode: 200,
	}
}
func (resp *Response) writeHeaderFunc(statusCode int) {
	resp.StatusCode = statusCode
	resp.writer.WriteHeader(statusCode)
}

func (resp *Response) writeFunc(data []byte) (int, error) {
	if !resp.Written {
		if resp.buffer != nil {
			resp.buffer.Check(resp)
		}
		resp.Written = true
	}

	length, err := resp.writer.Write(data)
	if err != nil {
		return 0, err
	}
	resp.ContentLength = resp.ContentLength + length
	if resp.buffer != nil {
		resp.buffer.Check(resp)
		resp.buffer.Write(data)
	}
	return length, nil
}

//Header http response writer header
func (resp *Response) Header() http.Header {
	return resp.writer.Header()
}

//WrapWriter wrap http response writer
func (resp *Response) WrapWriter(rw http.ResponseWriter) middleware.ResponseWriter {
	resp.writer = rw
	w := middleware.WrapResponseWriter(rw)
	f := w.Functions()
	f.WriteHeaderFunc = resp.writeHeaderFunc
	f.WriteFunc = resp.writeFunc
	return w
}

func (resp *Response) BuildBuffer(r *http.Request, v Validator) bool {
	return resp.BuildBufferWith(r, v, nil)
}

func (resp *Response) BuildBufferWith(r *http.Request, v Validator, writer io.Writer) bool {
	if resp.Written == true {
		return false
	}
	if resp.buffer != nil {
		return false
	}
	resp.buffer = NewResponseBuffer()
	resp.buffer.request = r
	if v != nil {
		resp.buffer.validator = v
	}
	if writer != nil {
		resp.buffer.writer = writer
	}
	return true
}

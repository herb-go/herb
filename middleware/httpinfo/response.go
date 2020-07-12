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
	async         bool
	Buffer        *bytes.Buffer
	pipe          Pipe
	locked        bool
}

//NewResponse create new response
func NewResponse() *Response {
	return &Response{
		StatusCode: 200,
		pipe:       DefaultPipe,
		header:     http.Header{},
		Buffer:     bytes.NewBuffer(nil),
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
	resp.flushHeader()
	resp.writer.WriteHeader(statusCode)
}

func (resp *Response) writeFunc(data []byte) (int, error) {
	resp.locked = true
	if !resp.Written {
		resp.pipe.Check()
		resp.Written = true
	}
	resp.locked = true
	length, err := resp.writer.Write(data)
	if err != nil {
		return 0, err
	}
	resp.ContentLength = resp.ContentLength + length
	resp.pipe.Check()
	resp.pipe.Write(data)
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

func (resp *Response) PipeDiscarded() bool {
	return resp.pipe.Discarded()
}

func (resp *Response) ReadAllBuffer() ([]byte, error) {
	if resp.Buffer == nil {
		return nil, nil
	}
	err := resp.pipe.Error()
	if err != nil {
		return nil, err
	}
	if resp.pipe.Discarded() {
		return nil, nil
	}
	return resp.Buffer.Bytes(), nil
}

func (resp *Response) Locked() bool {
	return resp.locked
}
func (resp *Response) UpdatePipe(p Pipe) bool {
	if resp.Locked() {
		return false
	}
	resp.pipe = p
	return true
}

// func (resp *Response) BuildBuffer(r *http.Request, v Validator) bool {
// 	if resp.Written == true {
// 		return false
// 	}
// 	if resp.buffer != nil {
// 		return false
// 	}
// 	resp.buffer = NewBuffer()
// 	resp.buffer.request = r
// 	if v != nil {
// 		resp.buffer.checker = v
// 	}
// 	return true
// }

// func (resp *Response) UpdateBufferWriter(writer io.Writer) bool {
// 	if !resp.bufferChangeable() {
// 		return false
// 	}
// 	if writer != nil {
// 		resp.buffer.writer = writer
// 	}
// 	return true
// }

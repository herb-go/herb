package responseinfo

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

type Info struct {
	StatusCode    int
	ContentLength int
	writer        http.ResponseWriter
}

func New() *Info {
	return &Info{
		StatusCode: 200,
	}
}
func (i *Info) writeHeaderFunc(statusCode int) {
	i.StatusCode = statusCode
	i.writer.WriteHeader(statusCode)
}

func (i *Info) writeFunc(data []byte) (int, error) {
	i.ContentLength = i.ContentLength + len(data)
	return i.writer.Write(data)
}
func (i *Info) Header() http.Header {
	return i.writer.Header()
}
func (i *Info) WrapWriter(rw http.ResponseWriter) middleware.ResponseWriter {
	i.writer = rw
	w := middleware.WrapResponseWriter(rw)
	f := w.Functions()
	f.WriteHeaderFunc = i.writeHeaderFunc
	f.WriteFunc = i.writeFunc
	return w
}

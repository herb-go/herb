package middleware

import (
	"bufio"
	"net"
	"net/http"
)

type WriterFunctions struct {
	WriteFunc       func([]byte) (int, error)
	HeaderFunc      func() http.Header
	WriteHeaderFunc func(statusCode int)
	FlushFunc       func()
	HijackFunc      func() (net.Conn, *bufio.ReadWriter, error)
}

type WrappedWriter struct {
	functions *WriterFunctions
}

func (w *WrappedWriter) Functions() *WriterFunctions {
	return w.functions
}

func (w *WrappedWriter) Write(b []byte) (int, error) {
	return w.functions.WriteFunc(b)
}

func (w *WrappedWriter) Header() http.Header {
	return w.functions.HeaderFunc()
}

func (w *WrappedWriter) WriteHeader(statusCode int) {
	w.functions.WriteHeaderFunc(statusCode)
}

func (w *WrappedWriter) Flush() {
	w.functions.FlushFunc()
}

func (w *WrappedWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.functions.HijackFunc()
}

type ResponseWriter interface {
	Functions() *WriterFunctions
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type Flusher interface {
	// Flush sends any buffered data to the client.
	Flush()
}
type Hijacker interface {
	Hijack() (net.Conn, *bufio.ReadWriter, error)
}

type WrappedResponseWriter struct {
	ResponseWriter
}

type WrappedResponseWriterHijacker struct {
	ResponseWriter
	Hijacker
}
type WrappedResponseWriterFlusher struct {
	ResponseWriter
	Flusher
}
type WrappedResponseWriterFlusherHijacker struct {
	ResponseWriter
	Flusher
	Hijacker
}

func WrapResponseWriter(rw http.ResponseWriter) ResponseWriter {
	var isFlusher bool
	var isHijacker bool
	w := &WrappedWriter{
		functions: &WriterFunctions{
			WriteFunc:       rw.Write,
			HeaderFunc:      rw.Header,
			WriteHeaderFunc: rw.WriteHeader,
		},
	}
	if f, ok := rw.(Flusher); ok {
		isFlusher = true
		w.functions.FlushFunc = f.Flush
	}
	if h, ok := rw.(Hijacker); ok {
		isHijacker = true
		w.functions.HijackFunc = h.Hijack
	}
	if isFlusher && isHijacker {
		return &WrappedResponseWriterFlusherHijacker{
			ResponseWriter: w,
			Flusher:        w,
			Hijacker:       w,
		}
	}
	if isFlusher {
		return &WrappedResponseWriterFlusher{
			ResponseWriter: w,
			Flusher:        w,
		}
	}
	if isHijacker {
		return &WrappedResponseWriterHijacker{
			ResponseWriter: w,
			Hijacker:       w,
		}
	}
	return &WrappedResponseWriter{
		ResponseWriter: w,
	}
}

package misc

import (
	"net/http"
	"strconv"
	"time"
)

//Writer http response writer interface.
type Writer interface {
	http.ResponseWriter
	http.Hijacker
}

//ElapsedTime add requset elapsed time to "Elapsed-Time" header of response.
//Elapsed time is time spent between middleware exetue and data wrote to response.
func ElapsedTime(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nw := elapsedTimeResponseWriter{
		Writer:    w.(Writer),
		Timestamp: time.Now().UnixNano(),
		written:   false,
	}
	next(&nw, r)

}

type elapsedTimeResponseWriter struct {
	Writer
	Timestamp int64
	written   bool
}

func (e *elapsedTimeResponseWriter) WriteHeader(status int) {
	if e.written == false {
		e.written = true
		e.Writer.Header().Set("Elapsed-Time", strconv.FormatInt(time.Now().UnixNano()-e.Timestamp, 10)+" ns")
	}
	e.Writer.WriteHeader(status)
}
func (e *elapsedTimeResponseWriter) Write(data []byte) (int, error) {
	if e.written == false {
		e.WriteHeader(http.StatusOK)
	}
	return e.Writer.Write(data)
}

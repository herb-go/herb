package misc

import (
	"net/http"
	"strconv"
	"time"
)

func ElapsedTime(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	nw := elapsedTimeResponseWriter{
		ResponseWriter: w,
		Timestamp:      time.Now().UnixNano(),
		written:        false,
	}
	next(&nw, r)

}

type elapsedTimeResponseWriter struct {
	http.ResponseWriter
	Timestamp int64
	written   bool
}

func (e *elapsedTimeResponseWriter) WriteHeader(status int) {
	if e.written == false {
		e.written = true
		e.ResponseWriter.Header().Set("Elapsed-Time", strconv.FormatInt(time.Now().UnixNano()-e.Timestamp, 10)+" ns")
	}
	e.ResponseWriter.WriteHeader(status)
}
func (e *elapsedTimeResponseWriter) Write(data []byte) (int, error) {
	if e.written == false {
		e.WriteHeader(http.StatusOK)
	}
	return e.ResponseWriter.Write(data)
}

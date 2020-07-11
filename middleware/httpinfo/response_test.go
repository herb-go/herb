package httpinfo

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/herb-go/herb/middleware"
)

var finishchan = make(chan int, 10)

func finishMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
	finishchan <- 1
}

var lastresp *Response

func lastrespMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	resp := NewResponse()
	next(resp.WrapWriter(w), r)
}

func echoAction(w http.ResponseWriter, r *http.Request) {
	for field := range r.Header {
		for k := range r.Header[field] {
			w.Header().Set(field, r.Header[field][k])
		}
	}
	status := r.Header.Get("status")
	if status != "" {
		statusCode, err := strconv.Atoi(status)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(statusCode)
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
func TestResponse(t *testing.T) {
	resp := NewResponse()
	if resp.StatusCode != 200 || resp.Written {
		t.Fatal(resp)
	}
	mux := http.NewServeMux()
	s := httptest.NewServer(mux)
	defer s.Close()
	mux.Handle("/test", middleware.New().Use(finishMiddleware, lastrespMiddleware).HandleFunc(echoAction))
	finishchan = make(chan int, 10)
	defer close(finishchan)
	req, err := http.NewRequest("POST", s.URL+"/test", bytes.NewBufferString("testcontent"))
	if err != nil {
		panic(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	<-finishchan
}

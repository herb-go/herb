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

func respMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	next(lastresp.WrapWriter(w), r)
}
func bufferMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	lastresp.BuildBuffer(r, nil)
	next(lastresp.WrapWriter(w), r)
}
func neverMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	lastresp.BuildBuffer(r, ValidatorNever)
	next(lastresp.WrapWriter(w), r)
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
	var data []byte
	resp := NewResponse()
	if resp.StatusCode != 200 || resp.Written {
		t.Fatal(resp)
	}
	mux := http.NewServeMux()
	s := httptest.NewServer(mux)
	defer s.Close()
	mux.Handle("/test", middleware.New().Use(finishMiddleware, respMiddleware).HandleFunc(echoAction))
	mux.Handle("/buffer", middleware.New().Use(finishMiddleware, bufferMiddleware).HandleFunc(echoAction))
	mux.Handle("/never", middleware.New().Use(finishMiddleware, neverMiddleware).HandleFunc(echoAction))

	finishchan = make(chan int, 10)
	defer close(finishchan)
	req, err := http.NewRequest("POST", s.URL+"/test", bytes.NewBufferString("testcontent"))
	if err != nil {
		panic(err)
	}
	req.Header.Set("status", "401")
	req.Header.Set("testfield", "testvalue")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	<-finishchan
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") || lastresp.BufferDiscarded() != true {
		t.Fatal(lastresp)
	}
	data, err = lastresp.ReadAllBuffer()
	if len(data) != 0 || err != nil {
		t.Fatal(data, err)
	}

	req, err = http.NewRequest("POST", s.URL+"/buffer", bytes.NewBufferString("testcontent"))
	if err != nil {
		panic(err)
	}
	req.Header.Set("status", "401")
	req.Header.Set("testfield", "testvalue")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	<-finishchan
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") || lastresp.BufferDiscarded() != false {
		t.Fatal(lastresp)
	}
	data, err = lastresp.ReadAllBuffer()
	if string(data) != "testcontent" || err != nil {
		t.Fatal(data, err)
	}
	req, err = http.NewRequest("POST", s.URL+"/never", bytes.NewBufferString("testcontent"))
	if err != nil {
		panic(err)
	}
	req.Header.Set("status", "401")
	req.Header.Set("testfield", "testvalue")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	<-finishchan
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") || lastresp.BufferDiscarded() != true {
		t.Fatal(lastresp)
	}
	data, err = lastresp.ReadAllBuffer()
	if len(data) != 0 || err != nil {
		t.Fatal(data, err)
	}
}

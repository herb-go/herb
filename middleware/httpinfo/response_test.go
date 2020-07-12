package httpinfo

import (
	"bytes"
	"errors"
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
	lastresp.UpdateController(NewBufferController(r, lastresp))
	next(lastresp.WrapWriter(w), r)
}
func neverMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	lastresp.UpdateController(
		NewBufferController(r, lastresp).
			WithChecker(ValidatorNever),
	)
	next(lastresp.WrapWriter(w), r)
}

var errtest = errors.New("errtest")

type errwriter struct {
}

func (w *errwriter) Write([]byte) (int, error) {
	return 0, errtest
}
func errwriterMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	lastresp.UpdateController(
		NewBufferController(r, lastresp).
			WithWriter(&errwriter{}),
	)
	next(lastresp.WrapWriter(w), r)
}

var validatorError = ValidatorFunc(func(*http.Request, *Response) (bool, error) {
	return true, errtest
})

func errvalidatorMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	lastresp.UpdateController(
		NewBufferController(r, lastresp).
			WithChecker(validatorError),
	)
	next(lastresp.WrapWriter(w), r)
}

var validatorNotWritten = ValidatorFunc(func(req *http.Request, resp *Response) (bool, error) {
	return !resp.Written, nil
})

func notwrittenvalidatorMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	lastresp.UpdateController(
		NewBufferController(r, lastresp).
			WithChecker(validatorNotWritten),
	)
	next(lastresp.WrapWriter(w), r)
}

var writerbuffer *bytes.Buffer

func writerbufferMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lastresp = NewResponse()
	writerbuffer = bytes.NewBuffer(nil)
	lastresp.UpdateController(
		NewBufferController(r, lastresp).
			WithChecker(ValidatorAlways).
			WithWriter(writerbuffer),
	)
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

func readAllBuffer() ([]byte, error) {
	err := lastresp.controller.Error()
	if err != nil {
		return nil, err
	}
	c := lastresp.controller
	b, ok := c.(*BufferController)
	if !ok {
		return nil, nil
	}
	return b.buffer.Bytes(), nil
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
	mux.Handle("/errwriter", middleware.New().Use(finishMiddleware, errwriterMiddleware).HandleFunc(echoAction))
	mux.Handle("/errvalidator", middleware.New().Use(finishMiddleware, errvalidatorMiddleware).HandleFunc(echoAction))
	mux.Handle("/notwritten", middleware.New().Use(finishMiddleware, notwrittenvalidatorMiddleware).HandleFunc(echoAction))
	mux.Handle("/writerbuffer", middleware.New().Use(finishMiddleware, writerbufferMiddleware).HandleFunc(echoAction))

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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
	if len(data) != 0 || err != nil {
		t.Fatal(data, err)
	}

	req, err = http.NewRequest("POST", s.URL+"/errwriter", bytes.NewBufferString("testcontent"))
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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
	if len(data) != 0 || err != errtest {
		t.Fatal(data, err)
	}

	req, err = http.NewRequest("POST", s.URL+"/errvalidator", bytes.NewBufferString("testcontent"))
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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
	if len(data) != 0 || err != errtest {
		t.Fatal(data, err)
	}

	req, err = http.NewRequest("POST", s.URL+"/notwritten", bytes.NewBufferString("testcontent"))
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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
	if len(data) != 0 || err != nil {
		t.Fatal(data, err)
	}

	req, err = http.NewRequest("POST", s.URL+"/writerbuffer", bytes.NewBufferString("testcontent"))
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
	if lastresp.StatusCode != 401 || lastresp.Header().Get("testfield") != "testvalue" || lastresp.ContentLength != len("testcontent") {
		t.Fatal(lastresp)
	}
	data, err = readAllBuffer()
	if len(data) != 0 || err != nil {
		t.Fatal(data, err)
	}
	data = writerbuffer.Bytes()
	if string(data) != "testcontent" {
		t.Fatal(data)
	}

}

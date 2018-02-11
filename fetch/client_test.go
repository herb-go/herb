package fetch

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("proxied"))
		if err != nil {
			panic(err)
		}
	}))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("notproxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer proxy.Close()
	clients := Clients{
		ProxyURL: proxy.URL,
	}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := clients.Fetch(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.BodyContent) != "proxied" {
		t.Error(string(resp.BodyContent))
	}
}

func TestNoProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("proxied"))
		if err != nil {
			panic(err)
		}
	}))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("notproxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer proxy.Close()
	clients := Clients{}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := clients.Fetch(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.BodyContent) != "notproxied" {
		t.Error(string(resp.BodyContent))
	}
}

func TestNilProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("proxied"))
		if err != nil {
			panic(err)
		}
	}))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("notproxied"))
		if err != nil {
			panic(err)
		}
	}))
	defer proxy.Close()
	var clients *Clients
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := clients.Fetch(req)
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.BodyContent) != "notproxied" {
		t.Error(string(resp.BodyContent))
	}
}

func TestUnmarshal(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/json", testJSONAction)
	mux.HandleFunc("/xml", testXMLAction)
	ts := httptest.NewServer(mux)
	defer ts.Close()
	server := Server{
		Host: ts.URL,
		Headers: http.Header{
			"test": []string{
				"testheader",
			},
		},
	}
	APIJSON := server.EndPoint("POST", "/json")
	APIXML := server.EndPoint("POST", "/xml")
	var clients *Clients
	req, err := APIJSON.NewJSONRequest(nil, testServerData{Name: "testname"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := clients.Fetch(req)
	if err != nil {
		t.Fatal(err)
	}
	result := &testServerResp{}
	err = resp.UnmarshalAsJSON(&result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Header != "testheader" {
		t.Error(result.Header)
	}
	if result.Name != "testname" {
		t.Error(result.Name)
	}

	req, err = APIXML.NewXMLRequest(nil, testServerData{Name: "testname"})
	if err != nil {
		t.Fatal(err)
	}
	resp, err = clients.Fetch(req)
	if err != nil {
		t.Fatal(err)
	}
	result = &testServerResp{}
	err = resp.UnmarshalAsXML(&result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Header != "testheader" {
		t.Error(result.Header)
	}
	if result.Name != "testname" {
		t.Error(result.Name)
	}

}

func TestResultAsErr(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(bytes.Repeat([]byte("0"), 1000))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := DefaultClients.Fetch(req)
	if err != nil {
		t.Fatal(err)
	}
	respinterface := interface{}(resp)
	resperr, ok := (respinterface).(error)
	if ok == false {
		t.Error(ok)
	}
	errmsg := resperr.Error()
	if len(errmsg) != ErrMsgLengthLimit {
		t.Error(len(errmsg))
	}
	code := GetErrorStatusCode(resp)
	if code != 200 {
		t.Error(code)
	}
	code = GetErrorStatusCode(*resp)
	if code != 200 {
		t.Error(code)
	}
	errTest := errors.New("test")
	code = GetErrorStatusCode(errTest)
	if code != 0 {
		t.Error(code)
	}
	apierr := resp.NewAPICodeErr(100)
	errmsg = apierr.Error()
	if len(errmsg) != ErrMsgLengthLimit {
		t.Error(len(errmsg))
	}

	respinterface = interface{}(apierr)
	resperr, ok = (respinterface).(error)
	if ok == false {
		t.Error(ok)
	}
	if GetAPIErrCode(apierr) != "100" {
		t.Error(GetAPIErrCode(apierr))
	}
	if CompareAPIErrCode(apierr, 100) != true {
		t.Error(apierr.Error())
	}
	if CompareAPIErrCode(apierr, "100") != true {
		t.Error(apierr.Error())
	}
	if GetAPIErrCode(*apierr) != "100" {
		t.Error(GetAPIErrCode(apierr))
	}
	if GetAPIErrCode(errTest) != "" {
		t.Error(GetAPIErrCode(errTest))
	}
	if CompareAPIErrCode(errTest, 100) == true {
		t.Error(errTest.Error())
	}

}

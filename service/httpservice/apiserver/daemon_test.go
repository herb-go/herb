package apiserver

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/herb-go/herb/server"
)

func newOption() *Option {
	return &Option{
		Server: Server{
			Name: "testServer",
			HTTPConfig: server.HTTPConfig{
				ListenerConfig: server.ListenerConfig{
					Net:  "tcp",
					Addr: "127.0.0.1:6789",
				},
			},
		},
		Channel: "test",
	}
}

func newOption2() *Option {
	return &Option{
		Server: Server{
			Name: "testServer2",
			HTTPConfig: server.HTTPConfig{
				ListenerConfig: server.ListenerConfig{
					Net:  "tcp",
					Addr: "127.0.0.1:6788",
				},
			},
		},
		Channel: "test",
	}
}
func TestRouter(t *testing.T) {
	var err error
	Reset()
	o := newOption()
	err = o.ApplyServer()
	if err != nil {
		t.Fatal(err)
	}
	as := o.server()
	as.CleanConfig()
	defer func() {
		Reset()
		as.CleanConfig()
	}()
	config := as.Config()
	as.Start("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/test"))
	})
	defer as.Stop("/test")
	as.Start("/test2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/test2"))
	})
	defer as.Stop("/test2")
	resp, err := http.Get("http://" + config.Addr + "/test")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "/test" {
		t.Fatal(string(bs))
	}
	resp, err = http.Get("http://" + config.Addr + "/test2")
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "/test2" {
		t.Fatal(string(bs))
	}
	resp, err = http.Get("http://" + config.Addr + "/notexist")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
}
func TestDaemon(t *testing.T) {
	var err error
	Reset()
	o := newOption()
	as := o.server()
	defer Reset()
	defer func() {
		as.CleanConfig()
	}()
	if as.runningCount != 0 {
		t.Fatal(as.runningCount)
	}
	err = as.Stop("test")
	if err == nil {
		t.Fatal(err)
	}
	err = as.Start("test", func(w http.ResponseWriter, r *http.Request) {})
	if err != nil {
		t.Fatal(err)
	}
	if as.runningCount != 1 {
		t.Fatal(as.runningCount)
	}
	err = as.Start("test", func(w http.ResponseWriter, r *http.Request) {})
	if err == nil {
		t.Fatal(err)
	}
	if as.runningCount != 1 {
		t.Fatal(as.runningCount)
	}
	err = as.Stop("test")
	if err != nil {
		t.Fatal(err)
	}
	if as.runningCount != 0 {
		t.Fatal(as.runningCount)
	}
	err = as.Start("test", func(w http.ResponseWriter, r *http.Request) {})
	if err != nil {
		t.Fatal(err)
	}
	err = as.Start("test2", func(w http.ResponseWriter, r *http.Request) {})
	if err != nil {
		t.Fatal(err)
	}
	if as.runningCount != 2 {
		t.Fatal(as.runningCount)
	}
	err = as.Stop("test")
	if err != nil {
		t.Fatal(err)
	}
	if as.runningCount != 1 {
		t.Fatal(as.runningCount)
	}
	err = as.Stop("test2")
	if err != nil {
		t.Fatal(err)
	}
	if as.runningCount != 0 {
		t.Fatal(as.runningCount)
	}
	err = as.Stop("test")
	if err == nil {
		t.Fatal(err)
	}
	if as.runningCount != 0 {
		t.Fatal(as.runningCount)
	}
}

func TestMutliOption(t *testing.T) {
	var err error
	Reset()
	o := newOption()
	err = o.ApplyServer()
	if err != nil {
		t.Fatal(err)
	}
	as := o.server()
	config := as.Config()

	defer func() {
		Reset()
		as.CleanConfig()
	}()
	o2 := newOption2()
	err = o2.ApplyServer()
	if err != nil {
		t.Fatal(err)
	}
	as2 := o2.server()
	config2 := as2.Config()
	defer func() {
		Reset()
		as2.CleanConfig()
	}()
	as.Start("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server1/test"))
	})
	defer as.Stop("/test")
	as.Start("/test1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server1/test1"))
	})
	defer as.Stop("/test1")
	as2.Start("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server2/test"))
	})
	defer as2.Stop("/test")
	as2.Start("/test2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server2/test2"))
	})
	defer as2.Stop("/test2")

	resp, err := http.Get("http://" + config.Addr + "/test")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "server1/test" {
		t.Fatal(string(bs))
	}

	resp, err = http.Get("http://" + config2.Addr + "/test")
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "server2/test" {
		t.Fatal(string(bs))
	}

	resp, err = http.Get("http://" + config.Addr + "/test1")
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "server1/test1" {
		t.Fatal(string(bs))
	}

	resp, err = http.Get("http://" + config2.Addr + "/test2")
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "server2/test2" {
		t.Fatal(string(bs))
	}

	resp, err = http.Get("http://" + config.Addr + "/test2")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}

	resp, err = http.Get("http://" + config2.Addr + "/test1")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
}

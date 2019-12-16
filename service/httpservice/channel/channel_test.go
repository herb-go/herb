package channel

import (
	"net/http"
	"testing"
)

func newTestChannel(prefix string) *Channel {
	c := NewChannel()
	c.ListenerConfig = testListener
	c.Path = prefix
	return c
}

func resetAll() {
	ResetServers()
	ResetConfigs()
	DefaultConfig = PresetDefaultConfig()
}
func TestDefaultChannel(t *testing.T) {
	// defer resetAll()
	// channel := NewChannel()
	// d := &url.URL{
	// 	Scheme: DefaultNet,
	// 	Host:   DefaulAddr,
	// }
	// if channel.Host() != d.String() {
	// 	t.Fatal(channel)
	// }
}

func TestChannel(t *testing.T) {
	result := []interface{}{}
	defer resetAll()
	channel := newTestChannel("/test")
	channel2 := newTestChannel("/test2/")
	channel.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result = append(result, "test")
		w.Write(nil)
	}))
	channel2.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result = append(result, r.URL.Path)
		w.Write(nil)
	}))
	err := channel.Start()
	if err != nil {
		t.Fatal(err)
	}
	path := "http://" + channel.Addr + "/test"
	resp, err := http.DefaultClient.Get(path)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	if len(result) != 1 || result[0].(string) != "test" {
		t.Fatal(result)
	}

	path2 := "http://" + channel.Addr + "/test2/test2path"
	resp, err = http.DefaultClient.Get(path2)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}

	err = channel2.Start()
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Get(path2)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal(resp)
	}
	if len(result) != 2 || result[1].(string) != "test2path" {
		t.Fatal(result)
	}
	err = channel2.Stop()
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Get(path2)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
	err = channel.Stop()
	if err != nil {
		t.Fatal(err)
	}
	_, err = http.DefaultClient.Get(path)
	if err == nil {
		t.Fatal(err)
	}

}
